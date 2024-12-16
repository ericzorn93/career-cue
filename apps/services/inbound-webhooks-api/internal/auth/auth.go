package auth

import (
	"go.uber.org/fx"
)

// NewAuthModule will register the authentication logic for
// the application
func NewAuthModule() fx.Option {
	authModule := fx.Module(
		"authModule",
		fx.Provide(
			fx.Annotate(NewAuthPublisher, fx.ResultTags(`name:"authPublisher"`)),
			fx.Annotate(NewService, fx.As(new(AuthServicePort)), fx.ResultTags(`name:"authService"`)),
			fx.Annotate(NewHandler, fx.As(new(AuthHandler)), fx.ResultTags(`name:"authHandler"`)),
		),
	)

	return authModule
}
