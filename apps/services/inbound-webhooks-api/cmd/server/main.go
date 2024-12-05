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

func NewServ() *serv {
	return &serv{}
}

func (s *serv) UserRegistered(ctx context.Context, req *pb.UserRegisteredRequest) (*commonv1.Empty, error) {
	log.Println("hit grpc endpoint")
	log.Println(req)
	return &commonv1.Empty{}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		server := grpc.NewServer()
		pb.RegisterInboundWebhooksAPIServer(server, NewServ())

		l, err := net.Listen("tcp", ":5000")
		if err != nil {
			log.Fatal(err)
		}

		if err := server.Serve(l); err != nil {
			log.Fatal(err)
		}
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
