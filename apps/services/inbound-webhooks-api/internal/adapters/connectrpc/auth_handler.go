package connectrpc

import (
	"context"

	"connectrpc.com/connect"

	"apps/services/inbound-webhooks-api/internal/app"
	boot "libs/backend/boot"
	"libs/backend/domain/user"
	commonv1 "libs/backend/proto-gen/go/common/v1"
	inboundwebhooksapiv1 "libs/backend/proto-gen/go/webhooks/inboundwebhooksapi/v1"
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
		user.NewUser(
			user.WithUserFirstName(req.Msg.FirstName),
			user.WithUserLastName(req.Msg.LastName),
			user.WithUserNickname(req.Msg.Nickname),
			user.WithUserUsername(req.Msg.Username),
			user.WithEmailAddress(req.Msg.EmailAddress),
			user.WithEmailAddressVerified(req.Msg.EmailAddressVerified),
			user.WithPhoneNumber(req.Msg.PhoneNumber),
			user.WithPhoneNumberVerified(req.Msg.PhoneNumberVerified),
			user.WithStrategy(req.Msg.Strategy),
			user.WithCommonID(req.Msg.CommonId),
			user.WithMetadata(make(map[string]any, 0)),
		),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
