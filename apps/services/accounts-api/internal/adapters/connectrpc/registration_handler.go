package connectrpc

import (
	"apps/services/accounts-api/internal/app"
	"context"
	"fmt"
	"libs/backend/boot"
	"libs/backend/domain/user/entities"
	"libs/backend/domain/user/valueobjects"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"

	"connectrpc.com/connect"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type RegistrationServiceHandler struct {
	accountsapiv1connect.UnimplementedRegistrationServiceHandler

	// Logger is the logger from the boot framework
	Logger boot.Logger

	// App is the application
	App app.App
}

// NewRegistrationHandler will return a pointer to the inbound webhooks API server
func NewRegistrationHandler(logger boot.Logger, app app.App) *RegistrationServiceHandler {
	return &RegistrationServiceHandler{
		Logger: logger,
		App:    app,
	}
}

// CreateAccount handles user creation and saves them in the database
func (r *RegistrationServiceHandler) CreateAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.CreateAccountRequest],
) (*connect.Response[accountsapiv1.CreateAcountResponse], error) {
	commonID, err := valueobjects.NewCommonIDFromString(req.Msg.CommonId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid common ID: %w", err))
	}

	// Convert to user domain type
	user := entities.NewUser(
		entities.WithCommonID(commonID),
		entities.WithEmailAddress(req.Msg.EmailAddress),
		entities.WithUserUsername(req.Msg.Username),
	)

	// Create a new user
	r.App.RegistrationService.RegisterUser(ctx, user)

	return connect.NewResponse(&accountsapiv1.CreateAcountResponse{IsSuccess: true}), nil
}
