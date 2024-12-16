package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	"log/slog"

	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
	boot "libs/boot/pkg"

	"go.uber.org/fx"
)

// InboundWebhooksAuthAPIServer handles all gRPC endpoints for inbound webhooks
type InboundWebhooksAuthAPIServer struct {
	pb.UnimplementedInboundWebhooksAuthAPIServer

	Service AuthServicePort
	Logger  boot.Logger
}

// New is the constructor for the inbound webhooks API
type NewHandlerParams struct {
	fx.In

	Service AuthServicePort `name:"authService"`
	Logger  boot.Logger
}

// NewHandler will return a pointer to the inbound webhooks API server
func NewHandler(params NewHandlerParams) *InboundWebhooksAuthAPIServer {
	return &InboundWebhooksAuthAPIServer{
		Service: params.Service,
		Logger:  params.Logger,
	}
}

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (s *InboundWebhooksAuthAPIServer) UserRegistered(
	ctx context.Context,
	req *pb.UserRegisteredRequest,
) (*commonv1.Empty, error) {
	if err := s.Service.SendUserRegistered(ctx, entities.NewUser()); err != nil {
		s.Logger.Info("cannot send user registered event", slog.Any("error", err))
		return nil, err
	}

	return &commonv1.Empty{}, nil
}
