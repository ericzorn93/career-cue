package services

import (
	"apps/services/accounts-worker/internal/app/ports"
	"context"
	"fmt"
	"libs/backend/auth/m2m"
	boot "libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	"libs/backend/httpauth"
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
	RegistrationServiceClient accountsapiv1connect.AccountServiceClient
	M2MClient                 m2m.M2MGenerator
}

// AccountServiceParams is a struct to hold the parameters for the AccountService
type AccountServiceParams struct {
	Logger          boot.Logger
	AccountConsumer boot.AMQPConsumer
	AccountsAPIURI  string
	M2MClient       m2m.M2MGenerator
}

// NewAccountService will construct the auth service
func NewAccountService(params AccountServiceParams) AccountService {
	registrationServiceClient := accountsapiv1connect.NewAccountServiceClient(http.DefaultClient, params.AccountsAPIURI)

	return AccountService{
		Logger:                    params.Logger,
		AccountConsumer:           params.AccountConsumer,
		RegistrationServiceClient: registrationServiceClient,
		M2MClient:                 params.M2MClient,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s AccountService) CreateAccount(ctx context.Context, user userEntities.User) error {
	// Obtain machine-to-machine-token
	m2mToken, err := s.M2MClient.GetToken()
	if err != nil {
		return fmt.Errorf("m2m token generation failure")
	}

	// Call the accounts-api to create the account
	req := connect.NewRequest(&accountsapiv1.CreateAccountRequest{
		Username:     user.Username,
		EmailAddress: user.EmailAddress.String(),
		CommonId:     user.CommonID.String(),
	})
	req.Header().Add(httpauth.AuthorizationHeaderKey, m2mToken.GetHeaderValue())
	account, err := s.RegistrationServiceClient.CreateAccount(ctx, req)

	if err != nil {
		s.Logger.Error("Cannot create account in Accounts API", slog.Any("error", err))
		return ports.ErrUserNotCreated
	}

	s.Logger.Info("Account created in Accounts API", slog.Any("isSuccess", account.Msg.GetIsSuccess()))
	return nil
}
