package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

// AuthService defines the core business logic
// for the auth service
type AuthService struct {
	Logger    *slog.Logger
	AuthQueue amqp.Queue
}

// AuthService params defines what the service can accept
type AuthServiceParams struct {
	fx.In

	Logger *slog.Logger

	// TODO: Add back
	// AuthQueue amqp.Queue `name:"authQueue"`
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
	time.Sleep(time.Millisecond * 800)
	s.Logger.Info("Calling the auth service SendUserRegistered event", slog.String("queueName", s.AuthQueue.Name))

	return nil
}
