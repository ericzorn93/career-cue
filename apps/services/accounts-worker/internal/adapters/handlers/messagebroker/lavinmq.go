package messagebroker

import (
	"apps/services/accounts-worker/internal/app"
	"context"
	"libs/backend/boot"
	"libs/backend/domain/user"
	accountseventsv1 "libs/backend/proto-gen/go/accounts/accountsevents/v1"
	"log/slog"

	"google.golang.org/protobuf/proto"
)

// LavinMQHandler handles all incoming events from LavinMQ
type LavinMQHandler struct {
	Logger   boot.Logger
	Consumer boot.AMQPConsumer
	App      app.App
}

// LavinMQHandler is the constructor for LavinMQHandler
func NewLavinMQHandler(logger boot.Logger, consumer boot.AMQPConsumer, app app.App) LavinMQHandler {
	return LavinMQHandler{
		Logger:   logger,
		Consumer: consumer,
		App:      app,
	}
}

// HandleUserRegisteredEvent handles the user registered event
func (h LavinMQHandler) HandleUserRegisteredEvent(ctx context.Context, queueName string) error {
	msgs, err := h.Consumer.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		h.Logger.Error("Cannot consume messages", slog.Any("error", err))
		return err
	}

	// Loop through all messages in queue
	for msg := range msgs {
		// Unmarshal the message to userRegisteredEvent
		var userRegisteredEvent accountseventsv1.UserRegistered
		proto.Unmarshal(msg.Body, &userRegisteredEvent)

		// Create User in Accounts API
		user := user.NewUser(
			user.WithUserFirstName(userRegisteredEvent.FirstName),
			user.WithUserLastName(userRegisteredEvent.LastName),
			user.WithUserNickname(userRegisteredEvent.Nickname),
			user.WithUserUsername(userRegisteredEvent.Username),
			user.WithEmailAddress(userRegisteredEvent.EmailAddress),
			user.WithEmailAddressVerified(userRegisteredEvent.EmailAddressVerified),
			user.WithPhoneNumber(userRegisteredEvent.PhoneNumber),
			user.WithPhoneNumberVerified(userRegisteredEvent.PhoneNumberVerified),
			user.WithStrategy(userRegisteredEvent.Strategy),
		)

		// Create Account in Accounts API
		go func() {
			if err := h.App.AccountService.CreateAccount(ctx, user); err != nil {
				h.Logger.Error("Cannot create account", slog.Any("error", err))
			}
		}()
	}

	return nil
}