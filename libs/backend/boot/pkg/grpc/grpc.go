package grpc

import (
	"context"
	"fmt"
	"libs/boot/pkg/logger"
	"log/slog"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/credentials"
)

// Handler is a type of callback used specifically for starting the gRPC handlers
type Handler func(context.Context, *http.ServeMux) error

// Options initializes how gRPC service gets started
type Options struct {
	Port                 uint64
	TransportCredentials []credentials.TransportCredentials
	GRPCHandlers         []Handler
	GatewayEnabled       bool
}

// startgRPCService will establish a TCP bound port and start the gRPC service
func StartgRPCService(ctx context.Context, serviceName string, logger logger.Logger, opts Options) error {
	// Check if the gRPC Gatway should exist
	if len(opts.GRPCHandlers) == 0 {
		logger.Info("No gRPC handlers present")
		return nil
	}

	// HTTP Server Mux
	mux := http.NewServeMux()

	// Register protobuf
	for _, grpcHandler := range opts.GRPCHandlers {
		if err := grpcHandler(ctx, mux); err != nil {
			logger.Error("Error occurred", "error", err)
			return err
		}
	}

	// Start Connect/gRPC Server
	logger.Info("Starting service on HTTP", slog.String("serviceName", serviceName))

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", opts.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		logger.Error("Error occurred", "error", err)
		return err
	}

	return nil
}
