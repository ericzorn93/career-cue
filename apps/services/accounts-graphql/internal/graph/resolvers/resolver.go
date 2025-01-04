package resolvers

import (
	"apps/services/accounts-graphql/internal/config"
	"libs/backend/boot"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	"net/http"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger            boot.Logger
	AccountsAPIClient accountsapiv1connect.AccountServiceClient
}

func NewResolver(logger boot.Logger, config config.Config) *Resolver {
	// Set up Accounts API Client
	accountsAPIClient := accountsapiv1connect.NewAccountServiceClient(
		http.DefaultClient,
		config.AccountsAPIURI,
	)

	return &Resolver{
		Logger:            logger,
		AccountsAPIClient: accountsAPIClient,
	}
}
