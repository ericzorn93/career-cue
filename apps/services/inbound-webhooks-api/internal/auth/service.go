package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AuthService defines the core business logic
// for the auth service
type AuthService struct {
	Logger    *slog.Logger
	AuthQueue amqp.Queue
}

// NewService is bound to the dependency injection framework and will initialize the
// auth service
func NewService(logger *slog.Logger, authQueue amqp.Queue) AuthServicePort {
	return AuthService{
		Logger: logger,
		// TODO: Fix this dependency injection error
		AuthQueue: authQueue,
	}
}

// SendUserRegistered is the concrete implementation of assigning the user registered event
// to the message broker
func (s AuthService) SendUserRegistered(ctx context.Context, user entities.User) error {
	time.Sleep(time.Millisecond * 800)
	s.Logger.Info("Calling the auth service SendUserRegistered event")

	return nil
}
