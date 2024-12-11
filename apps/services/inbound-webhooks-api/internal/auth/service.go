package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	"log/slog"
	"time"
)

// AuthService defines the core business logic
// for the auth service
type AuthService struct {
	Logger *slog.Logger
}

func NewService(logger *slog.Logger) AuthServicePort {
	return AuthService{
		Logger: logger,
	}
}

// SendUserRegistered is the concrete implementation of assigning the user registered event
// to the message broker
func (as AuthService) SendUserRegistered(ctx context.Context, user entities.User) error {
	time.Sleep(time.Millisecond * 800)
	as.Logger.Info("Calling the auth service SendUserRegistered event")
	return nil
}
