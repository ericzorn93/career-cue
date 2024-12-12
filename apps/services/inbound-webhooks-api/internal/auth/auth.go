package auth

import (
	"go.uber.org/fx"
)

// NewAuthModule will register the authentication logic for
// the application
func NewAuthModule() fx.Option {
	authModule := fx.Module(
		"auth",
		fx.Provide(NewHandler), // Handler must be public for main module
		fx.Provide(fx.Private, NewService),
		fx.Provide(fx.Private, EstablishQueues),
	)

	return authModule
}
