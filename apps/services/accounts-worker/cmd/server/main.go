package main

import (
	"apps/services/accounts-worker/internal/adapters/handlers/messagebroker"
	"apps/services/accounts-worker/internal/app"
	"apps/services/accounts-worker/internal/config"
	"apps/services/accounts-worker/internal/domain/services"
	"context"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	boot "libs/backend/boot"
	"libs/backend/eventing"
)

// setupHandler sets up the handler for the message broker
func setupHandler(logger boot.Logger, accountsAPIUri string, amqpConsumer boot.AMQPConsumer) messagebroker.LavinMQHandler {
	// Initialize the application
	accountService := services.NewAccountService(services.AccountServiceParams{
		Logger:         logger,
		AccountsAPIURI: accountsAPIUri,
	})
	application := app.NewApp(app.WithAccountService(accountService))

	// Initialize the message broker handler
	handler := messagebroker.NewLavinMQHandler(logger, amqpConsumer, application)

	return handler
}

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
	_ = []connect.HandlerOption{
		connect.WithInterceptors(validationInterceptor),
	}

	// Initialize the gRPC Options
	bootService := boot.
		NewBuildServiceBuilder().
		SetServiceName(config.ServiceName).
		SetLogger(logger).
		SetAMQPOptions(boot.AMQPOptions{
			ConnectionURI: config.AMQPUri,
			OnConnectionCallback: func(params boot.AMQPCallBackParams) error {
				params.Logger.Info("AMQP connected successfully")

				// Set up the auth event queues and exchanges
				authEventRegisterer := eventing.NewAuthEventSetup(params.Controller.Registerer, params.Logger)
				authEventRegisterer.
					CreateExchange().
					CreateQueue(config.UserRegistrationQueueName).
					BindQueues([]string{eventing.GetUserRegisteredRoutingKey()}).
					Complete()

				params.Logger.Info("Set up all AMQP queues and exchanges")

				return nil
			},
			Handlers: []boot.AMQPHandler{
				func(hp boot.AMQPHandlerParams) error {
					handler := setupHandler(hp.Logger, config.AccountsAPIUri, hp.AMQPController.Consumer)

					// Handle the user registered event
					if err := handler.HandleUserRegisteredEvent(ctx, config.UserRegistrationQueueName); err != nil {
						hp.Logger.Error("Cannot handle user registered event", slog.Any("error", err))
					}

					return nil
				},
			},
		}).
		SetConnectRPCOptions(boot.ConnectRPCOptions{
			Port: 3000,
			TransportCredentials: []credentials.TransportCredentials{
				insecure.NewCredentials(),
			},
			Handlers: []boot.ConnectRPCHandler{},
		}).
		SetBootCallbacks([]boot.BootCallback{
			func(params boot.BootCallbackParams) error {
				params.Logger.Info("Service booted successfully", slog.String("serviceName", config.ServiceName))
				return nil
			},
		}).
		Build()

	return bootService.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Cannot start service")
		os.Exit(1)
	}
}
