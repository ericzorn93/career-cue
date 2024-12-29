package connectrpc

import (
	"context"
	"log/slog"

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
	commonID, err := userValueObjects.NewCommonIDFromString(req.Msg.CommonId)
	if err != nil {
		h.Logger.Error("Cannot create common ID", slog.Any("error", err))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := h.Application.AuthService.RegisterUser(
		userEntities.NewUser(
			userEntities.WithUserFirstName(req.Msg.FirstName),
			userEntities.WithUserLastName(req.Msg.LastName),
			userEntities.WithUserNickname(req.Msg.Nickname),
			userEntities.WithUserUsername(req.Msg.Username),
			userEntities.WithEmailAddress(req.Msg.EmailAddress),
			userEntities.WithEmailAddressVerified(req.Msg.EmailAddressVerified),
			userEntities.WithPhoneNumber(req.Msg.PhoneNumber),
			userEntities.WithPhoneNumberVerified(req.Msg.PhoneNumberVerified),
			userEntities.WithStrategy(req.Msg.Strategy),
			userEntities.WithCommonID(commonID),
			userEntities.WithMetadata(make(map[string]any, 0)),
		),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&commonv1.Empty{}), nil
}
