package services

import (
	"apps/services/accounts-worker/internal/app/ports"
	"context"
	boot "libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
)

// AccountService handles generic interactions
type AccountService struct {
	Logger                    boot.Logger
	AccountConsumer           boot.AMQPConsumer
	RegistrationServiceClient accountsapiv1connect.RegistrationServiceClient
}

// AccountServiceParams is a struct to hold the parameters for the AccountService
type AccountServiceParams struct {
	Logger          boot.Logger
	AccountConsumer boot.AMQPConsumer
	AccountsAPIURI  string
}

// NewAccountService will construct the auth service
func NewAccountService(params AccountServiceParams) AccountService {
	registrationServiceClient := accountsapiv1connect.NewRegistrationServiceClient(http.DefaultClient, params.AccountsAPIURI)

	return AccountService{
		Logger:                    params.Logger,
		AccountConsumer:           params.AccountConsumer,
		RegistrationServiceClient: registrationServiceClient,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s AccountService) CreateAccount(ctx context.Context, user userEntities.User) error {
	// Call the accounts-api to create the account
	account, err := s.RegistrationServiceClient.CreateAccount(ctx, connect.NewRequest(&accountsapiv1.CreateAccountRequest{
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Nickname:             user.Nickname,
		Username:             user.Username,
		EmailAddress:         user.EmailAddress,
		EmailAddressVerified: user.EmailAddressVerified,
		PhoneNumber:          user.PhoneNumber,
		PhoneNumberVerified:  user.PhoneNumberVerified,
		CommonId:             user.CommonID.String(),
		Strategy:             user.Strategy,
	}))

	if err != nil {
		s.Logger.Error("Cannot create account in Accounts API", slog.Any("error", err))
		return ports.ErrUserNotCreated
	}

	s.Logger.Info("Account created in Accounts API", slog.Any("isSuccess", account.Msg.GetIsSuccess()))
	return nil
}
