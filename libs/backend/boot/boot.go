package boot

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	bootCallbacks     []BootCallback
}

// BootCallback are methods for when the service is booted
type BootCallbackParams struct {
	Logger Logger
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
	err := EstablishAMQPConnection(s, s.logger, opts)
	if err != nil {
		return err
	}

	return nil
}

// Start spins up the service
func (s *BootService) Start(ctx context.Context) error {
	s.logger.Info("Service started", slog.String("serviceName", s.name))

	// Handle Shutdown of the service
	go func(ctx context.Context, bs *BootService) {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
		<-exitCh
		bs.logger.Info("Shutting down the service")
		if err := bs.Stop(ctx); err != nil {
			bs.logger.Error("Trouble closing the boot service")
		}
	}(ctx, s)

	// Control Go Routines
	s.wg.Add(1)

	// Start the connection to AMQP broker
	if err := s.startAMQPBrokerConnection(s.amqpOptions); err != nil {
		s.logger.Error("Cannot establish connection to AMQP broker", slog.Any("error", err))
		os.Exit(1)
	}

	// Start the connectRPC Service
	go func() {
		defer s.wg.Done()

		if err := StartConnectRPCService(ctx, s.name, s.logger, s.amqpController, s.connectRPCOptions); err != nil {
			s.logger.Error("Cannot properly start connectRPC Service")
			os.Exit(1)
		}
	}()

	// Wait for services to start
	s.wg.Wait()

	// Execute boot callbacks after service starts
	for _, cb := range s.bootCallbacks {
		if err := cb(BootCallbackParams{
			Logger: s.logger,
		}); err != nil {
			s.logger.Error("failed to execute boot callback", slog.Any("error", err))
			os.Exit(1)
		}
	}

	return nil
}

// Close supports the io.Closer interface and will shutdown the service
// when the process is done or programatically
func (s *BootService) Close() error {
	if err := s.amqpController.Close(); err != nil {
		s.logger.Error("Trouble closing the AMQP connection")
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
