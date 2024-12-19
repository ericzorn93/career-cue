package connectrpc

import (
	"context"
	"fmt"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"
	"log/slog"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/credentials"
)

// HandlerParams will be passed to the handler registered
type HandlerParams struct {
	Context        context.Context
	Logger         logger.Logger
	Mux            *http.ServeMux
	AMQPController amqp.Controller
}

// Handler is a type of callback used specifically for starting the gRPC handlers
type Handler func(HandlerParams) error

// Options initializes how gRPC service gets started
type Options struct {
	Port                 uint64
	TransportCredentials []credentials.TransportCredentials
	Handlers             []Handler
	GatewayEnabled       bool
}

// startConnectRPCService will establish a TCP bound port and start the gRPC service
func StartConnectRPCService(ctx context.Context, serviceName string, logger logger.Logger, amqpController amqp.Controller, opts Options) error {
	// Check if the gRPC Gatway should exist
	if len(opts.Handlers) == 0 {
		logger.Info("No gRPC handlers present")
		return nil
	}

	// HTTP Server Mux
	mux := http.NewServeMux()

	// Register protobuf
	for _, grpcHandler := range opts.Handlers {
		if err := grpcHandler(HandlerParams{
			Context:        ctx,
			Logger:         logger,
			Mux:            mux,
			AMQPController: amqpController,
		}); err != nil {
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
