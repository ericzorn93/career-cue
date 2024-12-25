package services

import (
	"context"
	boot "libs/backend/boot"
	accountsapiv1 "libs/backend/proto-gen/go/accounts/accountsapi/v1"
	"libs/backend/proto-gen/go/accounts/accountsapi/v1/accountsapiv1connect"
	accountseventsv1 "libs/backend/proto-gen/go/accounts/accountsevents/v1"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

// AccountService handles generic interactions
type AccountService struct {
	Logger          boot.Logger
	AccountConsumer boot.AMQPConsumer
	AccountsAPIURI  string
}

// NewAccountService will construct the auth service
type AccountServiceParams struct {
	Logger          boot.Logger
	AccountConsumer boot.AMQPConsumer
	AccountsAPIURI  string
}

func NewAccountService(params AccountServiceParams) AccountService {
	return AccountService{
		Logger:          params.Logger,
		AccountConsumer: params.AccountConsumer,
		AccountsAPIURI:  params.AccountsAPIURI,
	}
}

// RegisterUser is an application interface method to handle user registration
// webhooks
func (s AccountService) PublishAccountCreated(queueName string) error {
	msgs, err := s.AccountConsumer.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		s.Logger.Error("Cannot consume messages", slog.Any("error", err))
		return err
	}

	for msg := range msgs {
		go temporaryHandleUserRegisteredEvent(s.Logger, msg, s.AccountsAPIURI)
	}

	return nil
}

// TODO: Remove after implementing the actual handler
func temporaryHandleUserRegisteredEvent(logger boot.Logger, msg amqp.Delivery, accountsAPIUri string) {
	var userRegisteredEvent accountseventsv1.UserRegistered
	proto.Unmarshal(msg.Body, &userRegisteredEvent)
	logger.Info("Received message", slog.String("msg", userRegisteredEvent.String()))

	// Call the accounts-api to create the account
	client := accountsapiv1connect.NewRegistrationServiceClient(http.DefaultClient, accountsAPIUri)
	client.CreateAccount(context.Background(), connect.NewRequest(&accountsapiv1.CreateAccountRequest{
		FirstName:            userRegisteredEvent.FirstName,
		LastName:             userRegisteredEvent.LastName,
		Nickname:             userRegisteredEvent.Nickname,
		Username:             userRegisteredEvent.Username,
		EmailAddress:         userRegisteredEvent.EmailAddress,
		EmailAddressVerified: userRegisteredEvent.EmailAddressVerified,
		PhoneNumber:          userRegisteredEvent.PhoneNumber,
		PhoneNumberVerified:  userRegisteredEvent.PhoneNumberVerified,
		Strategy:             userRegisteredEvent.Strategy,
	}))
}
