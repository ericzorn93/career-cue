package main

import (
	"apps/services/inbound-webhooks-api/internal/adapters/api"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Concurrently start the gRPC Service
	go func() {
		server := grpc.NewServer()
		api.RegisterServer(server)

		l, err := net.Listen("tcp", ":5000")
		if err != nil {
			log.Fatal(err)
		}

		// Enable server reflection
		reflection.Register(server)

		if err := server.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	// Handle HTTP Restful requests as a proxy to gRPC Service
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterInboundWebhooksAPIHandlerFromEndpoint(ctx, mux, ":5000", opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Inbound Webhooks Service")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
