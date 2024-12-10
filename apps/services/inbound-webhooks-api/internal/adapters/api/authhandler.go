package api

import (
	"context"
	"log"

	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
)

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (s *InboundWebhooksAPIServer) UserRegistered(
	_ context.Context,
	req *pb.UserRegisteredRequest,
) (*commonv1.Empty, error) {
	log.Println("hit grpc endpoint for user registered")

	return &commonv1.Empty{}, nil
}
