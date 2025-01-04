package connectrpc

import (
	"context"

	"connectrpc.com/connect"

	"apps/services/inbound-webhooks-api/internal/app"
	boot "libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
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
	// Parse CommonID
	commonID := userValueObjects.NewCommonIDFromString(req.Msg.CommonId)
	emailAddress := userValueObjects.NewEmailAddress(req.Msg.EmailAddress)

	if err := h.Application.AuthService.RegisterUser(
		userEntities.NewUser(
			userEntities.WithCommonID(commonID),
			userEntities.WithEmailAddress(emailAddress),
			userEntities.WithUserUsername(req.Msg.Username),
			userEntities.WithMetadata(make(map[string]any, 0)),
		),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
