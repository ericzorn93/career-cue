package connectrpc

import (
	"apps/services/inbound-webhooks-api/internal/application"
	"apps/services/inbound-webhooks-api/internal/domain"
	"context"

	commonv1 "libs/backend/proto-gen/go/common/v1"
	pb "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
	"libs/boot/pkg/logger"

	"connectrpc.com/connect"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type AuthHandler struct {
	Logger      logger.Logger
	AuthService application.AuthService
}

// NewAuthHandler will return a pointer to the inbound webhooks API server
func NewAuthHandler(logger logger.Logger, authService application.AuthService) *AuthHandler {
	return &AuthHandler{
		Logger:      logger,
		AuthService: authService,
	}
}

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (h *AuthHandler) UserRegistered(
	ctx context.Context,
	req *connect.Request[pb.UserRegisteredRequest],
) (*connect.Response[commonv1.Empty], error) {
	if err := h.AuthService.RegisterUser(domain.NewUser(domain.WithUserFirstName(req.Msg.FirstName), domain.WithUserLastName(req.Msg.LastName))); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
