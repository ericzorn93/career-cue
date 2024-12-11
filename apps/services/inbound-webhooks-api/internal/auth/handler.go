package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"

	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"

	"go.uber.org/fx"
)

// InboundWebhooksAuthAPIServer handles all gRPC endpoints for inbound webhooks
type InboundWebhooksAuthAPIServer struct {
	pb.UnimplementedInboundWebhooksAuthAPIServer
	service AuthServicePort
}

// New is the constructor for the inbound webhooks API
type NewHandlerParams struct {
	fx.In
	Service AuthServicePort
}

// NewHandler will return a pointer to the inbound webhooks API server
func NewHandler(params NewHandlerParams) *InboundWebhooksAuthAPIServer {
	return &InboundWebhooksAuthAPIServer{
		service: params.Service,
	}
}

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (s *InboundWebhooksAuthAPIServer) UserRegistered(
	ctx context.Context,
	req *pb.UserRegisteredRequest,
) (*commonv1.Empty, error) {
	s.service.SendUserRegistered(ctx, entities.NewUser())
	return &commonv1.Empty{}, nil
}
