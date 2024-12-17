package main

import (
	grpcAdapters "apps/services/inbound-webhooks-api/internal/adapters/grpc"
	"apps/services/inbound-webhooks-api/internal/application"
	"apps/services/inbound-webhooks-api/internal/eventing"
	"context"
	"log"
	"net/http"
	"os"

	"connectrpc.com/grpcreflect"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	inboundwebhooksapiv1connect "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1/inboundwebhooksapiv1connect"
	boot "libs/boot/pkg"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/grpc"
)

const serviceName = "inbound-webhooks-api"

func run() error {
	// Application Context
	ctx := context.Background()

	// Initialize the gRPC Options
	bootService, err := boot.NewBootService(
		boot.BootServiceParams{
			Name: serviceName,
			AMQPOptions: amqp.Options{
				ConnectionURI: "amqp://guest:guest@lavinmq:5672",
				OnConnectionCallback: func(params amqp.CallBackParams) error {
					params.Logger.Info("AMQP connected successfully")

					// Set Up Auth Events
					eventing.RegisterAuthEvents(eventing.NewAuthQueueParams{
						Log:     params.Logger,
						Channel: params.Channel,
					})

					params.Logger.Info("Set up all AMQP queues and exchanges")

					return nil
				},
			},
			BootCallbacks: []boot.BootCallback{
				func(params boot.BootCallbackParams) error {
					params.Logger.Info("Service booted successfully", "serviceName", serviceName)
					return nil
				},
			},
		},
	)
	if err != nil {
		return err
	}

	// Assign gRPC Options
	bootService.SetGRPCOptions(grpc.Options{
		Port: 3000,
		TransportCredentials: []credentials.TransportCredentials{
			insecure.NewCredentials(),
		},
		GRPCHandlers: []grpc.Handler{
			func(ctx context.Context, mux *http.ServeMux) error {
				logger := bootService.GetLogger()
				amqpController := bootService.GetAMQPController()

				authService := application.NewAuthServiceImpl(logger, amqpController.Publisher)
				authHandler := grpcAdapters.NewAuthHandler(logger, authService)
				path, handler := inboundwebhooksapiv1connect.NewInboundWebhooksAuthServiceHandler(authHandler)
				mux.Handle(path, handler)

				reflector := grpcreflect.NewStaticReflector(
					inboundwebhooksapiv1connect.InboundWebhooksAuthServiceName,
				)
				mux.Handle(grpcreflect.NewHandlerV1(reflector))
				mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
				return nil
			},
		},
	})

	return bootService.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Printf("Cannot start service %s", serviceName)
		os.Exit(1)
	}
}
