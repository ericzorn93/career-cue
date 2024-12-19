package main

import (
	connectrpcAdapters "apps/services/inbound-webhooks-api/internal/adapters/connectrpc"
	"apps/services/inbound-webhooks-api/internal/application"
	"apps/services/inbound-webhooks-api/internal/config"
	"apps/services/inbound-webhooks-api/internal/eventing"
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/grpcreflect"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	inboundwebhooksapiv1connect "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1/inboundwebhooksapiv1connect"
	boot "libs/boot/pkg"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/connectrpc"
	bootLogger "libs/boot/pkg/logger"
)

// serviceName is the name of the microservice
const serviceName = "inbound-webhooks-api"

func run() error {
	// Application Context
	ctx := context.Background()

	// Create logger
	logger := bootLogger.NewSlogger()

	// Construct config
	config, err := config.NewConfig()
	if err != nil {
		logger.Error("Trouble constructing config")
		os.Exit(1)
	}

	// Initialize the gRPC Options
	bootService := boot.
		NewBuildServiceBuilder().
		SetServiceName(serviceName).
		SetLogger(logger).
		SetAMQPOptions(amqp.Options{
			ConnectionURI: config.AMQPUrl,
			OnConnectionCallback: func(params amqp.CallBackParams) error {
				params.Logger.Info("AMQP connected successfully")

				// Set Up Auth Events
				eventing.RegisterAuthEvents(eventing.NewAuthQueueParams{
					Log:        params.Logger,
					Registerer: params.Controller.Registerer,
				})

				params.Logger.Info("Set up all AMQP queues and exchanges")

				return nil
			},
		}).
		SetConnectRPCOptions(connectrpc.Options{
			Port: config.RPCPort,
			TransportCredentials: []credentials.TransportCredentials{
				insecure.NewCredentials(),
			},
			Handlers: []connectrpc.Handler{
				func(params connectrpc.HandlerParams) error {
					if !params.AMQPController.IsConnected() {
						errMsg := "AMQP not conntected"
						logger.Error(errMsg)
						return errors.New(errMsg)
					}

					logger.Info("AMQP connected")
					authService := application.NewAuthServiceImpl(logger, params.AMQPController.Publisher)
					authHandler := connectrpcAdapters.NewAuthHandler(logger, authService)
					path, handler := inboundwebhooksapiv1connect.NewInboundWebhooksAuthServiceHandler(authHandler)
					params.Mux.Handle(path, handler)

					reflector := grpcreflect.NewStaticReflector(
						inboundwebhooksapiv1connect.InboundWebhooksAuthServiceName,
					)
					params.Mux.Handle(grpcreflect.NewHandlerV1(reflector))
					params.Mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
					return nil
				},
			},
		}).
		SetBootCallbacks([]boot.BootCallback{
			func(params boot.BootCallbackParams) error {
				params.Logger.Info("Service booted successfully", slog.String("serviceName", serviceName))
				return nil
			},
		}).
		Build()

	return bootService.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Cannot start service %s", serviceName)
		os.Exit(1)
	}
}
