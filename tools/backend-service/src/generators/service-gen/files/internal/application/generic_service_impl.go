package application

import (
	"libs/boot"
	"log/slog"
)

// GenericServiceImpl handles generic interactions
type GenericServiceImpl struct {
	Logger             boot.Logger
	AuthEventPublisher boot.AMQPPublisher
}

// NewAuthServiceImpl will construct the auth service
func NewAuthServiceImpl(logger boot.Logger, contoller boot.AMQPPublisher) GenericServiceImpl {
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
