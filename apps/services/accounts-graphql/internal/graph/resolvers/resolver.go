package resolvers

import (
	"apps/services/accounts-graphql/internal/config"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	"net/http"

	"connectrpc.com/connect"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AccountsAPIClient accountsapiv1connect.AccountServiceClient
}

func NewResolver(config config.Config) *Resolver {
	// Set up Accounts API Client
	accountsAPIClient := accountsapiv1connect.NewAccountServiceClient(
		http.DefaultClient,
		config.AccountsAPIURI,
		connect.WithGRPC(),
	)

	return &Resolver{
		AccountsAPIClient: accountsAPIClient,
	}
}
