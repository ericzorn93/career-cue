package auth

import (
	"context"

	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
	"libs/boot/pkg/logger"

	"connectrpc.com/connect"
)

// InboundWebhooksAuthAPIServer handles all gRPC endpoints for inbound webhooks
type InboundWebhooksAuthAPIServer struct {
	Logger logger.Logger
}

// NewHandler will return a pointer to the inbound webhooks API server
func NewHandler(logger logger.Logger) *InboundWebhooksAuthAPIServer {
	return &InboundWebhooksAuthAPIServer{
		Logger: logger,
	}
}

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (s *InboundWebhooksAuthAPIServer) UserRegistered(
	ctx context.Context,
	req *connect.Request[pb.UserRegisteredRequest],
) (*connect.Response[commonv1.Empty], error) {
	return connect.NewResponse(&commonv1.Empty{}), nil
}
