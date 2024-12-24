package connectrpc

import (
	"apps/services/inbound-webhooks-api/internal/application"
	"apps/services/inbound-webhooks-api/internal/domain"
	"context"

	boot "libs/backend/boot"
	commonv1 "libs/backend/proto-gen/go/common/v1"
	inboundwebhooksapiv1 "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"

	"connectrpc.com/connect"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type AuthHandler struct {
	inboundwebhooksapiv1.UnimplementedInboundWebhooksAuthAPIServer
	Logger      boot.Logger
	Application application.Application
}

// NewAuthHandler will return a pointer to the inbound webhooks API server
func NewAuthHandler(logger boot.Logger, application application.Application) *AuthHandler {
	return &AuthHandler{
		Logger:      logger,
		Application: application,
	}
}

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (h *AuthHandler) UserRegistered(
	ctx context.Context,
	req *connect.Request[inboundwebhooksapiv1.UserRegisteredRequest],
) (*connect.Response[commonv1.Empty], error) {
	if err := h.Application.AuthService.RegisterUser(
		domain.NewUser(
			domain.WithUserFirstName(req.Msg.FirstName),
			domain.WithUserLastName(req.Msg.LastName),
			domain.WithUserNickname(req.Msg.Nickname),
			domain.WithUserUsername(req.Msg.Username),
			domain.WithEmailAddress(req.Msg.EmailAddress),
			domain.WithEmailAddressVerified(req.Msg.EmailAddressVerified),
			domain.WithPhoneNumber(req.Msg.PhoneNumber),
			domain.WithPhoneNumberVerified(req.Msg.PhoneNumberVerified),
			domain.WithStrategy(req.Msg.Strategy),
			domain.WithMetadata(make(map[string]any, 0)),
		),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
