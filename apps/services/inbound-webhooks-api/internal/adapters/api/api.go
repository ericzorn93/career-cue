package api

import (
	pb "packages/proto-gen/go/webhooks/inboundwebhooksapi/v1"

	"google.golang.org/grpc"
)

// InboundWebhooksAPIServer handles all gRPC endpoints for inbound webhooks
type InboundWebhooksAPIServer struct {
	pb.UnimplementedInboundWebhooksAPIServer
}

// New is the constructor for the inbound webhooks API
func New() *InboundWebhooksAPIServer {
	return &InboundWebhooksAPIServer{}
}

// RegisterServer registers the inbound webhooks API struct with gRPC
func RegisterServer(server *grpc.Server) {
	pb.RegisterInboundWebhooksAPIServer(server, New())
}
