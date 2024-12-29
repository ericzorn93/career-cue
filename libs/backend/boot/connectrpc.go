package boot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
)

// ConnectRPCHandlerParams will be passed to the handler registered
type ConnectRPCHandlerParams struct {
	Context        context.Context
	Logger         Logger
	Mux            *http.ServeMux
	AMQPController AMQPController
	DB             *gorm.DB
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

	// Register Prometheus Metrics Handler
	mux.Handle("/metrics", promhttp.Handler())

	// Register protobuf
	for _, grpcHandler := range s.connectRPCOptions.Handlers {
		err := grpcHandler(ConnectRPCHandlerParams{
			Context:        ctx,
			Logger:         s.logger,
			Mux:            mux,
			AMQPController: s.amqpController,
			DB:             s.db,
		})

		if err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}
	}

	// Start Connect/gRPC Server
	s.logger.Info("Starting service on HTTP", slog.String("serviceName", s.name))

	// Create an error group to handle multiple goroutines running HTTP Service
	// Start the HTTP server
	go func() error {
		s.logger.Info("Service bound to port", slog.Uint64("port", s.connectRPCOptions.Port), slog.String("serviceName", s.name))

		if err := http.ListenAndServe(
			fmt.Sprintf(":%d", s.connectRPCOptions.Port),
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		); err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}

		return nil
	}()

	// Start the IPV6 bound HTTP server for Fly.io (production only) - Always run on port 5000
	go func() error {
		environment := os.Getenv("ENV")
		if environment != "production" && environment != "prod" {
			s.logger.Info("Not running in production on Fly.io, skipping IPV6 bound HTTP server")
			return nil
		}

		const flyPort uint64 = 8080
		s.logger.Info("Service bound to port in production", slog.Uint64("port", flyPort), slog.String("serviceName", s.name))

		if err := http.ListenAndServe(
			fmt.Sprintf("fly-local-6pn:%d", flyPort),
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		); err != nil {
			s.logger.Error("Error occurred", "error", err)
			return err
		}

		return nil
	}()

	return nil
}
