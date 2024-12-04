package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commonv1 "packages/proto-gen/go/common/v1"
	pb "packages/proto-gen/go/webhooks/inboundwebhooksapi/v1"
)

type serv struct {
	pb.UnimplementedInboundWebhooksAPIServer
}

func (s *serv) UserCreated(ctx context.Context, req *pb.UserCreatedRequest) (*commonv1.Empty, error) {
	log.Println("hit grpc endpoint")
	return &commonv1.Empty{}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		server := grpc.NewServer()
		pb.RegisterInboundWebhooksAPIServer(server, &serv{})

		l, _ := net.Listen("tcp", ":5000")
		server.Serve(l)
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterInboundWebhooksAPIHandlerFromEndpoint(ctx, mux, ":5000", opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Inbound Webhooks Service")
	http.ListenAndServe(":3000", mux)
}
