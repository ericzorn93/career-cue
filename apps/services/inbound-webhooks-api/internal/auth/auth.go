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
			fx.Annotate(NewAuthQueue, fx.ResultTags(`name:"authQueue"`)),
			fx.Annotate(NewService, fx.As(new(AuthServicePort))),
			NewHandler,
		),
	)

	return authModule
}
