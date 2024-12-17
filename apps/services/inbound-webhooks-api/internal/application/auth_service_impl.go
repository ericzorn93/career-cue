package application

import (
	"apps/services/inbound-webhooks-api/internal/constants"
	"apps/services/inbound-webhooks-api/internal/domain"
	"libs/boot/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AuthServiceImpl handles all application auth interactions
type AuthServiceImpl struct {
	Logger  logger.Logger
	Channel *amqp.Channel
}

// NewAuthServiceImpl will construct the auth service
func NewAuthServiceImpl(logger logger.Logger, channel *amqp.Channel) AuthServiceImpl {
	return AuthServiceImpl{
		Logger:  logger,
		Channel: channel,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s AuthServiceImpl) RegisterUser(user domain.User) {
	s.Logger.Info("Publishing userRegistered Event")
	s.Channel.Publish(constants.AuthExchangeName, constants.AuthQueueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello world"),
	})
}
