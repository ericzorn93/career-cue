package connectrpc

import (
	"apps/services/inbound-webhooks-api/internal/app"
	"apps/services/inbound-webhooks-api/internal/domain/entities"
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
	Application app.Application
}

// NewAuthHandler will return a pointer to the inbound webhooks API server
func NewAuthHandler(logger boot.Logger, application app.Application) *AuthHandler {
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
		entities.NewUser(
			entities.WithUserFirstName(req.Msg.FirstName),
			entities.WithUserLastName(req.Msg.LastName),
			entities.WithUserNickname(req.Msg.Nickname),
			entities.WithUserUsername(req.Msg.Username),
			entities.WithEmailAddress(req.Msg.EmailAddress),
			entities.WithEmailAddressVerified(req.Msg.EmailAddressVerified),
			entities.WithPhoneNumber(req.Msg.PhoneNumber),
			entities.WithPhoneNumberVerified(req.Msg.PhoneNumberVerified),
			entities.WithStrategy(req.Msg.Strategy),
			entities.WithMetadata(make(map[string]any, 0)),
		),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
