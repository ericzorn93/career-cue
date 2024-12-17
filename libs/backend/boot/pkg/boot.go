package boot

import (
	"context"
	"io"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/grpc"
	"libs/boot/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
)

// BootService primary struct that defines
// holding the data for all service modules
type BootService struct {
	io.Closer
	wg             *sync.WaitGroup
	name           string
	logger         logger.Logger
	gRPCOptions    grpc.Options
	amqpConnection *amqp091.Connection
	amqpChannel    *amqp091.Channel
	bootCallbacks  []BootCallback
}

// BootServiceParams are the incoming options
// for the boot service construction
type BootServiceParams struct {
	Name          string
	AMQPOptions   amqp.Options
	BootCallbacks []BootCallback
}

// BootCallback are methods for when the service is booted
type BootCallbackParams struct {
	Logger logger.Logger
}
type BootCallback func(BootCallbackParams) error

// NewBootService sets up constructor for the boot service
// without any functionality or options
func NewBootService(params BootServiceParams) (BootService, error) {
	// Setup logger
	logger := logger.NewSlogger()

	// Create BootService instance
	bs := BootService{
		wg:            new(sync.WaitGroup),
		name:          params.Name,
		logger:        logger,
		bootCallbacks: params.BootCallbacks,
	}

	// Start connection to AMQP Broker
	bs.startAMQPBrokerConnection(params.AMQPOptions)

	// Handle Shutdown of the service
	go func() {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
		<-exitCh

		logger.Info("Closing the LavinMQ connection")
		if err := bs.Stop(context.TODO()); err != nil {
			logger.Error("Trouble closing the boot service")
		}
	}()

	return bs, nil
}

// SetGRPCOptions will assign the gRPC options to the Boot Service
func (s *BootService) SetGRPCOptions(opts grpc.Options) {
	s.gRPCOptions = opts
}

// startAMQPBrokerConnection will assign the AMQP broker options to the Boot Service
// and happens on boot service construction/initialization
func (s *BootService) startAMQPBrokerConnection(opts amqp.Options) {
	// Establish connection to AMQP broker
	amqpConnection, amqpChannel, err := amqp.EstablishAMQPConnection(s.logger, opts)
	if err != nil {
		s.logger.Error("Cannot establish connection to AMQP broker", slog.Any("error", err))
		return
	}

	// Assign Connections
	s.amqpConnection = amqpConnection
	s.amqpChannel = amqpChannel
}

// Start spins up the service
func (s BootService) Start(ctx context.Context) error {
	s.logger.Info("Service started", "serviceName", s.name)

	// Control Go Routines
	s.wg.Add(1)

	// Start the gRPC Service
	go func() {
		defer s.wg.Done()
		if err := grpc.StartgRPCService(ctx, s.name, s.logger, s.gRPCOptions); err != nil {
			s.logger.Error("cannot properly start gRPC Service")
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
func (s BootService) Close() error {
	return nil
}

// Stop proxies the call to the io.Closer method
// and will spin down the service
func (s BootService) Stop(_ context.Context) error {
	s.logger.Info("Service stopped", "serviceName", s.name)
	return s.Close()
}

// GetServiceName returns the name of the service
func (s BootService) GetServiceName() string {
	return s.name
}

// GetLogger returns the logger to the caller
func (s BootService) GetLogger() logger.Logger {
	return s.logger
}

// GetAMQPChannel will return the the channel associated with the AMQP broker connection
func (s BootService) GetAMQPChannel() *amqp091.Channel {
	return s.amqpChannel
}
