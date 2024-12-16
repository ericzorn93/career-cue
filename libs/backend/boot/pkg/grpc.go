package boot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/credentials"
)

// GRPCHandler is a type of callback used specifically for starting the gRPC handlers
type GRPCHandler func(context.Context, *http.ServeMux) error

// GRPCOptions initializes how gRPC service gets started
type GRPCOptions struct {
	Port                 uint64
	TransportCredentials []credentials.TransportCredentials
	GRPCHandlers         []GRPCHandler
	GatewayEnabled       bool
}

// startgRPCService will establish a TCP bound port and start the gRPC service
func (s BootService) startgRPCService(ctx context.Context) error {
	// Check if the gRPC Gatway should exist
	if len(s.gRPCOptions.GRPCHandlers) == 0 {
		s.log.Info("No gRPC handlers present")
		return nil
	}

	// HTTP Server Mux
	mux := http.NewServeMux()

	// Register protobuf
	for _, grpcHandler := range s.gRPCOptions.GRPCHandlers {
		if err := grpcHandler(ctx, mux); err != nil {
			s.log.Error("Error occurred", "error", err)
			return err
		}
	}

	// Start Connect/gRPC Server
	s.log.Info("Starting service on HTTP", slog.String("serviceName", s.name))

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.gRPCOptions.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		s.log.Error("Error occurred", "error", err)
		return err
	}

	return nil
}
