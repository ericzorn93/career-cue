package main

import (
	"apps/services/inbound-webhooks-api/internal/auth"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/grpcreflect"
	"go.uber.org/fx"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	inboundwebhooksapiv1connect "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1/inboundwebhooksapiv1connect"
	boot "libs/boot/pkg"
)

const serviceName = "inbound-webhooks-api"

func run() error {
	service := fx.Module(
		serviceName,
		auth.NewAuthModule(),
		fx.Provide(fx.Annotate(
			func(log boot.Logger, authHandler auth.AuthHandler) boot.BootServiceParams {
				return boot.BootServiceParams{
					Name: serviceName,
					GRPCOptions: boot.GRPCOptions{
						Port: 3000,
						TransportCredentials: []credentials.TransportCredentials{
							insecure.NewCredentials(),
						},
						GRPCHandlers: []boot.GRPCHandler{
							func(ctx context.Context, mux *http.ServeMux) error {
								path, handler := inboundwebhooksapiv1connect.NewInboundWebhooksAuthServiceHandler(authHandler)
								mux.Handle(path, handler)
								return nil
							},
							func(ctx context.Context, mux *http.ServeMux) error {
								reflector := grpcreflect.NewStaticReflector(
									inboundwebhooksapiv1connect.InboundWebhooksAuthServiceName,
								)
								mux.Handle(grpcreflect.NewHandlerV1(reflector))
								mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
								return nil
							},
						},
					},
					LavinMQOptions: boot.LavinMQOptions{
						ConnectionURI: "amqp://guest:guest@lavinmq:5672",
						OnConnectionCallback: func() error {
							log.Info("LavinMQ connected successfully")
							return nil
						},
					},
					BootCallbacks: []boot.BootCallback{
						func() error {
							log.Info("Service booted successfully", "serviceName", serviceName)
							return nil
						},
					},
				}
			},
			fx.ParamTags(``, `name:"authHandler"`),
		)),
		boot.NewBootServiceModule(),
	)

	app := fx.New(
		service,
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
	)
	app.Run()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("Cannot start service %s", serviceName)
		os.Exit(1)
	}
}
