package main

import (
	connectrpcAdapters "apps/services/inbound-webhooks-api/internal/adapters/connectrpc"
	"apps/services/inbound-webhooks-api/internal/app"
	"apps/services/inbound-webhooks-api/internal/app/usecases"
	"apps/services/inbound-webhooks-api/internal/config"
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	boot "libs/backend/boot"
	"libs/backend/eventing"
	inboundwebhooksapiv1connect "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1/inboundwebhooksapiv1connect"
)

// serviceName is the name of the microservice
const serviceName = "inbound-webhooks-api"

func run() error {
	// Application Context
	ctx := context.Background()

	// Create logger
	logger := boot.NewSlogger()

	// Construct config
	config, err := config.NewConfig()
	if err != nil {
		logger.Error("Trouble constructing config")
		os.Exit(1)
	}

	// Connect Interceptors
	validationInterceptor, err := validate.NewInterceptor()
	if err != nil {
		logger.Error("Cannot set up validation interceptor", slog.Any("error", err))
		return err
	}
	options := []connect.HandlerOption{
		connect.WithInterceptors(validationInterceptor),
	}

	// Initialize the gRPC Options
	bootService := boot.
		NewBuildServiceBuilder().
		SetServiceName(serviceName).
		SetLogger(logger).
		SetAMQPOptions(boot.AMQPOptions{
			ConnectionURI: config.AMQPUrl,
			OnConnectionCallback: func(params boot.AMQPCallBackParams) error {
				params.Logger.Info("AMQP connected successfully")

				// Set up the auth event queues and exchanges
				authEventRegisterer := eventing.NewAuthEventSetup(params.Controller.Registerer, params.Logger)
				authEventRegisterer.
					CreateExchange().
					CreateDeadletter().
					Complete()

				params.Logger.Info("Set up all AMQP queues and exchanges")

				return nil
			},
		}).
		SetConnectRPCOptions(boot.ConnectRPCOptions{
			Port: 3000,
			TransportCredentials: []credentials.TransportCredentials{
				insecure.NewCredentials(),
			},
			Handlers: []boot.ConnectRPCHandler{
				func(params boot.ConnectRPCHandlerParams) error {
					if !params.AMQPController.IsConnected() {
						errMsg := "AMQP not conntected"
						logger.Error(errMsg)
						return errors.New(errMsg)
					}

					// Construct application
					authService := usecases.NewAuthService(logger, params.AMQPController.Publisher)
					application := app.NewApplication(authService)
					authHandler := connectrpcAdapters.NewAuthHandler(logger, application)

					// Set up auth service routes and handlers
					path, handler := inboundwebhooksapiv1connect.NewInboundWebhooksAuthServiceHandler(
						authHandler,
						options...,
					)
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
