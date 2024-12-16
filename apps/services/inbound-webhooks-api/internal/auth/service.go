package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	boot "libs/boot/pkg"

	"go.uber.org/fx"
)

// AuthService defines the core business logic
// for the auth service
type AuthService struct {
	Logger boot.Logger
}

// AuthService params defines what the service can accept
type AuthServiceParams struct {
	fx.In

	Logger boot.Logger
}

// NewService is bound to the dependency injection framework and will initialize the
// auth service
func NewService(params AuthServiceParams) AuthServicePort {
	return AuthService{
		Logger: params.Logger,
	}
}

// SendUserRegistered is the concrete implementation of assigning the user registered event
// to the message broker
func (s AuthService) SendUserRegistered(ctx context.Context, user entities.User) error {
	s.Logger.Info("Calling the auth service SendUserRegistered event")

	return nil
}
