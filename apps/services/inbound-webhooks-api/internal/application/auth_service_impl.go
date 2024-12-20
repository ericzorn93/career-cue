package application

import (
	"apps/services/inbound-webhooks-api/internal/domain"
	"libs/backend/eventing"
	accountseventsv1 "libs/backend/proto-gen/go/accounts/accountsevents/v1"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
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
func (s AuthServiceImpl) RegisterUser(user domain.User) error {
	s.Logger.Info("Publishing userRegistered Event")

	metadata := make(map[string]*anypb.Any)
	for key, val := range user.Metadata {
		if convertedVal, err := anypb.New(structpb.NewStringValue(val.(string))); err != nil {
			s.Logger.Debug("Cannnot convert value in struct")
		} else {
			metadata[key] = convertedVal
		}
	}

	// Create and send event
	userRegisteredEvent := &accountseventsv1.UserRegistered{
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Nickname:             user.Nickname,
		Username:             user.Username,
		EmailAddress:         user.EmailAddress,
		EmailAddressVerified: user.EmailAddressVerified,
		PhoneNumber:          user.PhoneNumber,
		PhoneNumberVerified:  user.PhoneNumberVerified,
		Strategy:             user.Strategy,
		UserMetadata:         metadata,
	}
	b, err := proto.Marshal(userRegisteredEvent)
	if err != nil {
		s.Logger.Error("Cannot marshal user registered event")
		return err
	}

	s.AuthEventPublisher.Publish(eventing.AuthExchange, eventing.GetUserRegisteredRoutingKey(), false, false, amqp091.Publishing{
		ContentType: "application/x-protobuf",
		Body:        b,
	})

	return nil
}
