package boot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
)

// ConnectRPCHandlerParams will be passed to the handler registered
type ConnectRPCHandlerParams struct {
	Context        context.Context
	Logger         Logger
	Mux            *http.ServeMux
	AMQPController AMQPController
}

// ConnectRPCHandler is a type of callback used specifically for starting the gRPC handlers
type ConnectRPCHandler func(ConnectRPCHandlerParams) error

// ConnectRPCOptions initializes how gRPC service gets started
type ConnectRPCOptions struct {
	Port                 uint64
	TransportCredentials []credentials.TransportCredentials
	Handlers             []ConnectRPCHandler
	GatewayEnabled       bool
}

// StartConnectRPCService will establish a TCP bound port and start the gRPC service
func (s *BootService) StartConnectRPCService(ctx context.Context) error {
	// Check if the gRPC Gatway should exist
	if len(s.connectRPCOptions.Handlers) == 0 {
		s.logger.Info("No gRPC handlers present")
		return nil
	}

	// HTTP Server Mux
	mux := http.NewServeMux()

	// Register protobuf
	for _, grpcHandler := range s.connectRPCOptions.Handlers {
		err := grpcHandler(ConnectRPCHandlerParams{
			Context:        ctx,
			Logger:         s.logger,
			Mux:            mux,
			AMQPController: s.amqpController,
		})

		if err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}
	}

	// Start Connect/gRPC Server
	s.logger.Info("Starting service on HTTP", slog.String("serviceName", s.name))

	// Create an error group to handle multiple goroutines running HTTP Service
	egroup := errgroup.Group{}
	egroup.SetLimit(2)

	// Start the HTTP server
	egroup.Go(func() error {
		if err := http.ListenAndServe(
			fmt.Sprintf(":%d", s.connectRPCOptions.Port),
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		); err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}

		return nil
	})

	// Start the IPV6 bound HTTP server for Fly.io (production only)
	egroup.Go(func() error {
		// Get IPV6 address from environment
		flyPrivateIP := os.Getenv("FLY_PRIVATE_IP")

		if flyPrivateIP == "" {
			s.logger.Info("No Fly.io private IP found, running in Development mode")
			return nil
		}

		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", flyPrivateIP, s.connectRPCOptions.Port),
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		); err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}

		return nil
	})

	return egroup.Wait()
}
