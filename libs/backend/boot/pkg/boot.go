package boot

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
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
		fx.Provide(NewBootService),
		fx.Invoke(func(lc fx.Lifecycle, bs BootService, log *slog.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("Service started", "serviceName", bs.name)
					return bs.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					log.Info("Service stopped", "serviceName", bs.name)
					return bs.Stop(ctx)
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
	name          string
	log           *slog.Logger
	gRPCOptions   GRPCOptions
	bootCallbacks []BootCallback
}

// BootServiceParams are the incoming options
// for the boot service construction
type BootServiceParams struct {
	Name          string
	GRPCOptions   GRPCOptions
	bootCallbacks []BootCallback
}

type BootCallback func(context.Context) error
type GRPCHandler func(context.Context, *grpc.Server) error
type GatewayHandler func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

// GRPCOptions initializes how gRPC service gets started
type GRPCOptions struct {
	Port                 uint64
	GatewayPort          uint64
	TransportCredentials []credentials.TransportCredentials
	ReflectionEnabled    bool
	GRPCHandlers         []GRPCHandler
	GatewayEnabled       bool
	GatewayHandlers      []GatewayHandler
}

// NewBootService sets up constructor for the boot service
// without any functionality or options
func NewBootService(params BootServiceParams, log *slog.Logger) BootService {
	return BootService{
		name:          params.Name,
		log:           log,
		gRPCOptions:   params.GRPCOptions,
		bootCallbacks: params.bootCallbacks,
	}
}

// Close supports the io.Closer interface and will shutdown the service
// when the process is done or programatically
func (s *BootService) Close() error {
	return nil
}

// Start spins up the service
func (s BootService) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(10)

	// Spins up gRPC Service
	eg.Go(func() error {
		// Check if the gRPC Gatway should exist
		if len(s.gRPCOptions.GRPCHandlers) == 0 {
			s.log.Info("No gRPC handlers present")
			return nil
		}

		// gRPC Server Listen
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.gRPCOptions.Port))
		if err != nil {
			return err
		}

		// Make Credentials
		creds := make([]grpc.ServerOption, 0, len(s.gRPCOptions.TransportCredentials))
		for _, cred := range s.gRPCOptions.TransportCredentials {
			creds = append(creds, grpc.Creds(cred))
		}

		// Establish gRPC Server
		server := grpc.NewServer(creds...)

		// Register protobuf
		for _, grpcHandler := range s.gRPCOptions.GRPCHandlers {
			if err := grpcHandler(ctx, server); err != nil {
				s.log.Error("Error occurred", "error", err)
				return err
			}
		}

		// Enable server reflection
		if s.gRPCOptions.ReflectionEnabled {
			reflection.Register(server)
		}

		// Start gRPC Server
		if err := server.Serve(l); err != nil {
			s.log.Error("Error occurred", "error", err)
			return err
		}

		return nil
	})

	// If gateways present, spin up gRPC gatway proxies
	eg.Go(func() error {
		if len(s.gRPCOptions.GatewayHandlers) == 0 {
			return nil
		}

		mux := runtime.NewServeMux()

		dialOptions := make([]grpc.DialOption, 0, len(s.gRPCOptions.TransportCredentials))
		for _, cred := range s.gRPCOptions.TransportCredentials {
			dialOptions = append(dialOptions, grpc.WithTransportCredentials(cred))
		}

		for _, gatewayHandler := range s.gRPCOptions.GatewayHandlers {
			err := gatewayHandler(ctx, mux, fmt.Sprintf(":%d", s.gRPCOptions.Port), dialOptions)
			if err != nil {
				s.log.Error("Error occurred", "error", err)
				return err
			}
		}

		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.gRPCOptions.GatewayPort), mux); err != nil {
			s.log.Error("Error occurred", "error", err)
			return err
		}

		return nil
	})

	return eg.Wait()
}

// Stop proxies the call to the io.Closer method
// and will spin down the service
func (s BootService) Stop(ctx context.Context) error {
	s.log.Info("Shutting down service")
	return s.Close()
}
