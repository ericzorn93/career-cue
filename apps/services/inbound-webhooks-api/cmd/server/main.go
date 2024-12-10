package main

import (
	"apps/services/inbound-webhooks-api/internal/adapters/api"
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
	boot "libs/boot/pkg"
)

const serviceName = "inbound-webhooks-api"

func main() {
	service := fx.Module(
		serviceName,
		fx.Provide(func() boot.BootServiceParams {
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

	// // Concurrently start the gRPC Service
	// go func() {
	// 	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	// 	api.RegisterServer(server)

	// 	l, err := net.Listen("tcp", ":5000")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// Enable server reflection
	// 	reflection.Register(server)

	// 	if err := server.Serve(l); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// Handle HTTP Restful requests as a proxy to gRPC Service
	// mux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// err := pb.RegisterInboundWebhooksAPIHandlerFromEndpoint(ctx, mux, ":5000", opts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Starting Inbound Webhooks Service")
	// if err := http.ListenAndServe(":3000", mux); err != nil {
	// 	log.Fatal(err)
	// }
}
