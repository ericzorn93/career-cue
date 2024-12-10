package boot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// GRPCHandler is a type of callback used specifically for starting the gRPC handlers
type GRPCHandler func(context.Context, *grpc.Server) error

// GatewayHandler is called as a callback when the gRPC proxy server is established
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

// startgRPCService will establish a TCP bound port and start the gRPC service
func (s BootService) startgRPCService(ctx context.Context) error {
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
}

// startGRPCGateway will concurrently set up a proxy server to the underlying gRPC service
func (s BootService) startGRPCGateway(ctx context.Context) error {
	s.log.Info("starting gRPC gateway")

	if len(s.gRPCOptions.GatewayHandlers) == 0 {
		return nil
	}

	// Create server mux for HTTP requests
	mux := runtime.NewServeMux()

	// Convert transport credentials to dial options
	dialOptions := make([]grpc.DialOption, 0, len(s.gRPCOptions.TransportCredentials))
	for _, cred := range s.gRPCOptions.TransportCredentials {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(cred))
	}

	// Attach gateway protobug handlers
	for _, gatewayHandler := range s.gRPCOptions.GatewayHandlers {
		err := gatewayHandler(ctx, mux, fmt.Sprintf(":%d", s.gRPCOptions.Port), dialOptions)
		if err != nil {
			s.log.Error("Error occurred", "error", err)
			return err
		}
	}

	// Start HTTP server for REST proxy requests
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.gRPCOptions.GatewayPort), mux); err != nil {
		s.log.Error("Error occurred", "error", err)
		os.Exit(1)
	}

	return nil
}
