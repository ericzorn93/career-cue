package main

import (
	"apps/services/accounts-worker/internal/config"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	boot "libs/backend/boot"
	"libs/backend/eventing"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	accountseventsv1 "libs/backend/proto-gen/go/accounts/accountsevents/v1"
)

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

				// Set up all AMQP queues and exchanges
				if err := eventing.CreateUserRegisterationAuthEventInfrastructure(eventing.RegisterAuthParams{
					Registerer: params.Controller.Registerer,
					Log:        params.Logger,
					RoutingKey: eventing.GetUserRegisteredRoutingKey(),
					QueueName:  config.UserRegistrationQueueName,
				}); err != nil {
					params.Logger.Error("Cannot set up auth events", slog.Any("error", err))
				}

				params.Logger.Info("Set up all AMQP queues and exchanges")

				return nil
			},
			Handlers: []boot.AMQPHandler{
				func(hp boot.AMQPHandlerParams) error {
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
						go temporaryHandleUserRegisteredEvent(hp.Logger, msg, config.AccountsAPIUri)
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

// TODO: Remove after implementing the actual handler
func temporaryHandleUserRegisteredEvent(logger boot.Logger, msg amqp.Delivery, accountsAPIUri string) {
	var userRegisteredEvent accountseventsv1.UserRegistered
	proto.Unmarshal(msg.Body, &userRegisteredEvent)
	logger.Info("Received message", slog.String("msg", userRegisteredEvent.String()))

	// Call the accounts-api to create the account
	client := accountsapiv1connect.NewRegistrationServiceClient(http.DefaultClient, accountsAPIUri)
	client.CreateAccount(context.Background(), connect.NewRequest(&accountsapiv1.CreateAccountRequest{
		FirstName:            userRegisteredEvent.FirstName,
		LastName:             userRegisteredEvent.LastName,
		Nickname:             userRegisteredEvent.Nickname,
		Username:             userRegisteredEvent.Username,
		EmailAddress:         userRegisteredEvent.EmailAddress,
		EmailAddressVerified: userRegisteredEvent.EmailAddressVerified,
		PhoneNumber:          userRegisteredEvent.PhoneNumber,
		PhoneNumberVerified:  userRegisteredEvent.PhoneNumberVerified,
		Strategy:             userRegisteredEvent.Strategy,
	}))
}
