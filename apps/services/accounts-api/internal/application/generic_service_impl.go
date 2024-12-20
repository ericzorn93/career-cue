package application

import (
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"
	"log/slog"
)

// GenericServiceImpl handles generic interactions
type GenericServiceImpl struct {
	Logger             logger.Logger
	AuthEventPublisher amqp.Publisher
}

// NewAuthServiceImpl will construct the auth service
func NewAuthServiceImpl(logger logger.Logger, contoller amqp.Publisher) GenericServiceImpl {
	return GenericServiceImpl{
		Logger:             logger,
		AuthEventPublisher: contoller,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s GenericServiceImpl) RandomMethod(greeting string) error {
	s.Logger.Info("Publishing handling random greeting", slog.Any("greeting", greeting))
	return nil
}
