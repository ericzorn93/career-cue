package connectrpc

import (
	"apps/services/accounts-api/internal/app"
	"context"
	"fmt"
	"libs/backend/boot"
	"libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	accountsDomain "libs/backend/proto-gen/go/accounts/domain"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler handles all gRPC endpoints for inbound webhooks
type AccountServiceHandler struct {
	accountsapiv1connect.UnimplementedAccountServiceHandler

	// Logger is the logger from the boot framework
	Logger boot.Logger

	// App is the application
	App app.App
}

// NewRegistrationHandler will return a pointer to the inbound webhooks API server
func NewRegistrationHandler(logger boot.Logger, app app.App) *AccountServiceHandler {
	return &AccountServiceHandler{
		Logger: logger,
		App:    app,
	}
}

// CreateAccount handles user creation and saves them in the database
func (r *AccountServiceHandler) CreateAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.CreateAccountRequest],
) (*connect.Response[accountsapiv1.CreateAcountResponse], error) {
	commonID := userValueObjects.NewCommonIDFromString(req.Msg.CommonId)
	emailAddress := userValueObjects.NewEmailAddress(req.Msg.EmailAddress)

	// Convert to user domain type
	user := entities.NewUser(
		entities.WithCommonID(commonID),
		entities.WithEmailAddress(emailAddress),
		entities.WithUserUsername(req.Msg.Username),
	)

	// Create a new user
	r.App.RegistrationService.RegisterUser(ctx, user)

	return connect.NewResponse(&accountsapiv1.CreateAcountResponse{IsSuccess: true}), nil
}

// GetAccount handles user creation and saves them in the database
func (r *AccountServiceHandler) GetAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.GetAccountRequest],
) (*connect.Response[accountsapiv1.GetAccountResponse], error) {
	// Parse the commonID vs emailAddress, depending on which one is passed
	var commonID userValueObjects.CommonID
	var emailAddress userValueObjects.EmailAddress

	if req.Msg.CommonId != nil {
		commonID = userValueObjects.NewCommonIDFromString(*req.Msg.CommonId)
	}

	if req.Msg.EmailAddress != nil {
		emailAddress = userValueObjects.NewEmailAddress(*req.Msg.EmailAddress)
	}

	// Get the user from the database
	user, err := r.App.RegistrationService.GetUser(ctx, commonID, emailAddress)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found: %w", err))
	}

	// Convert to proto types
	account := &accountsDomain.Account{
		CommonId:     user.CommonID.String(),
		EmailAddress: user.EmailAddress.String(),
		Username:     user.Username,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}
	resp := &accountsapiv1.GetAccountResponse{Account: account}

	return connect.NewResponse(resp), nil
}

// DeleteAccount handles account deletion
func (r *AccountServiceHandler) DeleteAccount(
	ctx context.Context,
	req *connect.Request[accountsapiv1.DeleteAccountRequest],
) (*connect.Response[accountsapiv1.DeleteAccountResponse], error) {
	// Parse incoming request data
	parsedCommonID := userValueObjects.NewCommonIDFromString(req.Msg.CommonId)
	hardDelete := req.Msg.HardDelete

	// Perform deletion logic
	deletedAt, err := r.App.RegistrationService.DeleteUser(ctx, parsedCommonID, hardDelete)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&accountsapiv1.DeleteAccountResponse{DeletedAt: timestamppb.New(deletedAt)}), nil
}
