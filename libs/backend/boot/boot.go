package boot

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gorm.io/gorm"
)

// BootServiceBuilder is a builder struct for the BootService Instance
type BootServiceBuilder struct {
	bootService *BootService
}

// NewBuilServiceBuilder is a constructor to eventually build a boot service
func NewBuildServiceBuilder() *BootServiceBuilder {
	wg := &sync.WaitGroup{}
	bootService := &BootService{wg: wg}
	return &BootServiceBuilder{bootService: bootService}
}

// SetServiceName sets the name of the microservice on the BootService
func (bsb *BootServiceBuilder) SetServiceName(serviceName string) *BootServiceBuilder {
	bsb.bootService.name = serviceName
	return bsb
}

// SetLogger sets the logger on the BootService
func (bsb *BootServiceBuilder) SetLogger(logger Logger) *BootServiceBuilder {
	bsb.bootService.logger = logger
	return bsb
}

// SetConnectRPCOptions sets the connectRPC Options for connection on the BootService
func (bsb *BootServiceBuilder) SetConnectRPCOptions(connectRPCOptions ConnectRPCOptions) *BootServiceBuilder {
	bsb.bootService.connectRPCOptions = connectRPCOptions
	return bsb
}

// SetAMQPOptions sets the AMQP broker options for connection on the BootService
func (bsb *BootServiceBuilder) SetAMQPOptions(amqpOptions AMQPOptions) *BootServiceBuilder {
	bsb.bootService.amqpOptions = amqpOptions
	return bsb
}

// SetDBOptions sets the DB options for connection on the BootService
func (bsb *BootServiceBuilder) SetDBOptions(dbOptions DBOptions) *BootServiceBuilder {
	bsb.bootService.dbOptions = dbOptions
	return bsb
}

// SetBootCallbacks sets the boot callback for connection on the BootService
func (bsb *BootServiceBuilder) SetBootCallbacks(bootCallbacks []BootCallback) *BootServiceBuilder {
	bsb.bootService.bootCallbacks = bootCallbacks
	return bsb
}

// Build will eventually build the entire boot service struct in a complete format
func (bsb *BootServiceBuilder) Build() BootService {
	return *bsb.bootService
}

// BootService primary struct that defines
// holding the data for all service modules
type BootService struct {
	io.Closer
	wg                *sync.WaitGroup
	name              string
	logger            Logger
	connectRPCOptions ConnectRPCOptions
	amqpOptions       AMQPOptions
	amqpController    AMQPController
	dbOptions         DBOptions
	localDB           *sql.DB
	db                *gorm.DB
	bootCallbacks     []BootCallback
}

// BootCallback are methods for when the service is booted
type BootCallbackParams struct {
	Logger Logger
	DB     *gorm.DB
}
type BootCallback func(BootCallbackParams) error

// startAMQPBrokerConnection will assign the AMQP broker options to the Boot Service
// and happens on boot service construction/initialization
func (s *BootService) startAMQPBrokerConnection(opts AMQPOptions) error {
	// Check if options are empty
	if opts.IsZero() {
		s.logger.Warn("AMQP will not be used with empty config")
		return nil
	}

	// Establish connection to AMQP broker
	err := s.EstablishAMQPConnection(opts)
	if err != nil {
		s.logger.Error("Cannot establish connection to AMQP broker", slog.Any("error", err))
		return err
	}

	return nil
}

// startDBConnection will assign the DB options to the Boot Service
// and start the connection to the DB
func (s *BootService) startDBConnection(opts DBOptions) error {
	// Check if options are empty
	if opts.IsZero() {
		s.logger.Warn("DB will not be used with empty config")
		return nil
	}

	// Establish connection to DB
	s.logger.Info("Establishing connection to DB")
	return s.InitializeDB()
}

// Start spins up the service
func (s *BootService) Start(ctx context.Context) error {
	forever := make(chan struct{})
	defer close(forever)
	s.logger.Info("Service started", slog.String("serviceName", s.name))

	// Handle Shutdown of the service
	go func(ctx context.Context, bs *BootService) {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
		<-exitCh
		bs.logger.Info("Shutting down the service")
		if err := bs.Stop(ctx); err != nil {
			bs.logger.Error("Trouble closing the boot service", slog.String("serviceName", bs.name))
		}
		os.Exit(0)
	}(ctx, s)

	// Start the connection to AMQP broker
	if err := s.startAMQPBrokerConnection(s.amqpOptions); err != nil {
		s.logger.Error("Cannot establish connection to AMQP broker", slog.Any("error", err))
		os.Exit(1)
	}

	// Start the connection to CockroachDB
	if err := s.startDBConnection(s.dbOptions); err != nil {
		s.logger.Error("Cannot establish connection to DB", slog.Any("error", err))
		os.Exit(1)
	}

	// Start the connectRPC Service
	if err := s.StartConnectRPCService(ctx); err != nil {
		s.logger.Error("Cannot properly start connectRPC Service")
		os.Exit(1)
	}

	// Execute boot callbacks after service starts
	for _, cb := range s.bootCallbacks {
		if err := cb(BootCallbackParams{
			Logger: s.logger,
			DB:     s.db,
		}); err != nil {
			s.logger.Error("failed to execute boot callback", slog.Any("error", err))
			os.Exit(1)
		}
	}

	<-forever
	return nil
}

// Close supports the io.Closer interface and will shutdown the service
// when the process is done or programatically
func (s *BootService) Close() error {
	// Close the AMQP connection
	if err := s.amqpController.Close(); err != nil {
		s.logger.Error("Trouble closing the AMQP connection")
	}

	// Close the CockroachDB connection
	s.logger.Info("Closing DB connection")
	if err := s.localDB.Close(); err != nil {
		s.logger.Error("Trouble closing the DB connection")
	}

	s.logger.Info("Service stopped", "serviceName", s.name)
	return nil
}

// Stop proxies the call to the io.Closer method
// and will spin down the service
func (s BootService) Stop(_ context.Context) error {
	return s.Close()
}

// GetServiceName returns the name of the service
func (s BootService) GetServiceName() string {
	return s.name
}

// GetLogger returns the logger to the caller
func (s BootService) GetLogger() Logger {
	return s.logger
}

// GetAMQPController will return the the channel associated with the AMQP broker connection
func (s BootService) GetAMQPController() AMQPController {
	return s.amqpController
}
