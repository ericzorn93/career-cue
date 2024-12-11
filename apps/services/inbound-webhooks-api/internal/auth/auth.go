package auth

import (
	"go.uber.org/fx"
)

// NewAuthModule will register the authentication logic for
// the application
func NewAuthModule() fx.Option {
	authModule := fx.Module(
		"auth",
		fx.Provide(NewHandler),
		fx.Provide(NewService),
	)

	return authModule
}
