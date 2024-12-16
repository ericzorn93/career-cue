package boot

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"

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
		NewLoggerModule(),
		NewLavinMQModule(),
		fx.Provide(NewBootService),
		fx.Invoke(func(lc fx.Lifecycle, bs BootService, log Logger) {
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
	name           string
	log            Logger
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
func NewBootService(params BootServiceParams, log Logger) BootService {
	return BootService{
		wg:             new(sync.WaitGroup),
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
	// Control Go Routines
	s.wg.Add(1)

	// Start the gRPC Service
	go func() {
		defer s.wg.Done()
		if err := s.startgRPCService(ctx); err != nil {
			s.log.Error("cannot properly start gRPC Service")
			os.Exit(1)
		}
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
