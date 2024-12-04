package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "packages/proto-gen/go/webhooks/inboundwebhooksapi/v1"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterInboundWebhooksAPIHandlerFromEndpoint(ctx, mux, "localhost:3000", opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Inbound Webhooks Service")
	http.ListenAndServe(":3000", mux)
}
