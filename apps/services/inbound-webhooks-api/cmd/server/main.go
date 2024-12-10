package main

import (
	"apps/services/inbound-webhooks-api/internal/adapters/api"
	"context"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
	boot "libs/boot/pkg"

	amqp "github.com/rabbitmq/amqp091-go"
)

const serviceName = "inbound-webhooks-api"

func main() {
	service := fx.Module(
		serviceName,
		fx.Provide(func(log *slog.Logger) boot.BootServiceParams {
			return boot.BootServiceParams{
				Name: serviceName,
				GRPCOptions: boot.GRPCOptions{
					Port:        3000,
					GatewayPort: 5000,
					TransportCredentials: []credentials.TransportCredentials{
						insecure.NewCredentials(),
					},
					ReflectionEnabled: true,
					GRPCHandlers: []boot.GRPCHandler{
						func(ctx context.Context, srv *grpc.Server) error {
							api.RegisterServer(srv)
							return nil
						},
					},
					GatewayHandlers: []boot.GatewayHandler{
						func(ctx context.Context, mux *runtime.ServeMux, port string, dialOpts []grpc.DialOption) error {
							return pb.RegisterInboundWebhooksAPIHandlerFromEndpoint(ctx, mux, port, dialOpts)
						},
					},
				},
				LavinMQOptions: boot.LavinMQOptions{
					ConnectionURI: "amqp://guest:guest@lavinmq:5672",
					OnConnectionCallback: func(_ *amqp.Connection) error {
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
		}),
		boot.NewBootServiceModule(),
	)

	app := fx.New(
		service,
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
	)
	app.Run()
}
