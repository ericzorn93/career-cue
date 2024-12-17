package application

import (
	"apps/services/inbound-webhooks-api/internal/constants"
	"apps/services/inbound-webhooks-api/internal/domain"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

// AuthServiceImpl handles all application auth interactions
type AuthServiceImpl struct {
	Logger             logger.Logger
	AuthEventPublisher amqp.Publisher
}

// NewAuthServiceImpl will construct the auth service
func NewAuthServiceImpl(logger logger.Logger, contoller amqp.Publisher) AuthServiceImpl {
	return AuthServiceImpl{
		Logger:             logger,
		AuthEventPublisher: contoller,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s AuthServiceImpl) RegisterUser(user domain.User) {
	s.Logger.Info("Publishing userRegistered Event")
	s.AuthEventPublisher.Publish(constants.AuthExchangeName, constants.AuthQueueName, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello world"),
	})
}
