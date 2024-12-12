package boot

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

const (
	bootServiceModuleName = "bootService"
)

// NewBootServiceModule helps define the dependency injection that
// can easily be connected within any microservice system
func NewBootServiceModule() fx.Option {
	module := fx.Module(
		bootServiceModuleName,
		fx.Provide(func() *slog.Logger {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: false,
			}))
			return logger
		}),
		fx.Provide(func(lc fx.Lifecycle, bsParams BootServiceParams, log *slog.Logger) *amqp.Connection {
			// Validate the conneciton URI to LavinMQ
			if bsParams.LavinMQOptions.ConnectionURI == "" || !strings.HasPrefix(bsParams.LavinMQOptions.ConnectionURI, "amqp://") {
				log.Error("Cannot connect to LavinMQ with empty connection string")
				return nil
			}

			conn, err := amqp.DialConfig(bsParams.LavinMQOptions.ConnectionURI, amqp.Config{
				Heartbeat: time.Second * 30,
			})
			if err != nil {
				log.Error("Cannot connect to LavinMQ", slog.Any("error", err))
				os.Exit(1)
			}

			// Stop the connection on close
			lc.Append(fx.StopHook(func() error {
				return conn.Close()
			}))

			// Start LavinMQ Connection
			if err := bsParams.LavinMQOptions.OnConnectionCallback(); err != nil {
				log.Error("LavinMQ connection callback failed", slog.Any("error", err))
				os.Exit(1)
			}

			return conn
		}),
		fx.Provide(NewBootService),
		fx.Invoke(func(lc fx.Lifecycle, bs BootService, log *slog.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					log.Info("Service started", "serviceName", bs.name)
					go bs.Start(context.Background())
					return nil
				},
				OnStop: func(_ context.Context) error {
					log.Info("Service stopped", "serviceName", bs.name)
					go bs.Stop(context.Background())
					return nil
				},
			})
		}),
	)

	return module
}

// BootService primary struct that defines
// holding the data for all service modules
type BootService struct {
	io.Closer
	wg             *sync.WaitGroup
	mu             *sync.Mutex
	name           string
	log            *slog.Logger
	gRPCOptions    GRPCOptions
	lavinMQOptions LavinMQOptions
	bootCallbacks  []BootCallback
}

// BootServiceParams are the incoming options
// for the boot service construction
type BootServiceParams struct {
	Name           string
	GRPCOptions    GRPCOptions
	LavinMQOptions LavinMQOptions
	BootCallbacks  []BootCallback
}

// BootCallback are methods for when the service is booted
type BootCallback func() error

// NewBootService sets up constructor for the boot service
// without any functionality or options
func NewBootService(params BootServiceParams, log *slog.Logger) BootService {
	return BootService{
		wg:             new(sync.WaitGroup),
		mu:             new(sync.Mutex),
		name:           params.Name,
		log:            log,
		gRPCOptions:    params.GRPCOptions,
		lavinMQOptions: params.LavinMQOptions,
		bootCallbacks:  params.BootCallbacks,
	}
}

// Close supports the io.Closer interface and will shutdown the service
// when the process is done or programatically
func (s BootService) Close() error {
	return nil
}

// Start spins up the service
func (s BootService) Start(ctx context.Context) error {
	s.wg.Add(3)

	// Start the gRPC Service
	go func() {
		defer s.wg.Done()
		if err := s.startgRPCService(ctx); err != nil {
			s.log.ErrorContext(ctx, "cannot properly start gRPC Service")
			os.Exit(1)
		}
	}()

	// If gateway is present, spin up gRPC gateway proxies
	go func() {
		defer s.wg.Done()
		if err := s.startGRPCGateway(ctx); err != nil {
			s.log.ErrorContext(ctx, "cannot properly start gRPC Service")
			os.Exit(1)
		}
	}()

	// Establish a LavinMQ connection and pass the connection back to the caller
	// if the caller prefers to use LavinMQ as a producer and/or consumer in their
	// service
	go func() {
		defer s.wg.Done()

	}()

	// Wait for services to start
	s.wg.Wait()

	// Execute boot callbacks after service starts
	for _, cb := range s.bootCallbacks {
		if err := cb(); err != nil {
			s.log.Error("failed to execute boot callback", slog.Any("error", err))
			os.Exit(1)
		}
	}

	return nil
}

// Stop proxies the call to the io.Closer method
// and will spin down the service
func (s BootService) Stop(_ context.Context) error {
	s.log.Info("Shutting down service")
	return s.Close()
}
