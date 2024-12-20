package main

import (
	"apps/services/accounts-worker/internal/config"
	"context"
	"log"
	"log/slog"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"libs/backend/eventing"
	accountseventsv1 "libs/backend/proto-gen/go/accounts/accountsevents/v1"
	boot "libs/boot/pkg"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/connectrpc"
	bootLogger "libs/boot/pkg/logger"
)

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
		SetAMQPOptions(amqp.Options{
			ConnectionURI: config.AMQPUrl,
			OnConnectionCallback: func(params amqp.CallBackParams) error {
				params.Logger.Info("AMQP connected successfully")

				// Set up all AMQP queues and exchanges
				eventing.CreateUserRegisterationAuthEventInfrastructure(eventing.RegisterAuthParams{
					Registerer: params.Controller.Registerer,
					Log:        params.Logger,
					RoutingKey: eventing.GetUserRegisteredRoutingKey(),
					QueueName:  config.UserRegistrationQueueName,
				})

				params.Logger.Info("Set up all AMQP queues and exchanges")

				return nil
			},
			Handlers: []amqp.Handler{
				func(hp amqp.HandlerParams) error {
					msgs, err := hp.AMQPController.Consumer.Consume(
						config.UserRegistrationQueueName, // queue
						"",                               // consumer
						true,                             // auto-ack
						false,                            // exclusive
						false,                            // no-local
						false,                            // no-wait
						nil,                              // args
					)
					if err != nil {
						hp.Logger.Error("Cannot consume messages", slog.Any("error", err))
						return err
					}

					for msg := range msgs {
						var userRegisteredEvent accountseventsv1.UserRegistered
						proto.Unmarshal(msg.Body, &userRegisteredEvent)
						hp.Logger.Info("Received message", slog.Any("msg", userRegisteredEvent))
					}

					return nil
				},
			},
		}).
		SetConnectRPCOptions(connectrpc.Options{
			Port: 3000,
			TransportCredentials: []credentials.TransportCredentials{
				insecure.NewCredentials(),
			},
			Handlers: []connectrpc.Handler{},
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
