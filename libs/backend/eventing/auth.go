package eventing

import (
	boot "libs/backend/boot"
	"log/slog"
)

// Event Producer/Consumer
const (
	AuthExchange = "authExchange"
	AuthDomain   = "auth"
)

// Event Names
var (
	EventNameUserRegistered EventName = EventName(GetEventName(AuthDomain, "userRegistered"))
)

// Routing Keys
var (
	defaultAuthRoutingKey = GetRoutingKeyPrefix(AuthDomain) + ".*"
)

// GetUserRegisteredRoutingKey returns the routing key for user registered event
func GetUserRegisteredRoutingKey() string {
	return EventNameUserRegistered.String()
}

// RegisterAuthParams are params for the auth queue constructor
type RegisterAuthParams struct {
	Registerer boot.AMQPRegisterer
	Log        boot.Logger
	QueueName  string
	RoutingKey string
}

// CreateUserRegisterationAuthEventInfrastructure constructs Auth Queue from AMQP Channel
func CreateUserRegisterationAuthEventInfrastructure(params RegisterAuthParams) error {
	// Initialize Auth Exchange - topic
	err := params.Registerer.ExchangeDeclare(AuthExchange, "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		return err
	}

	// Create a queue for the service
	if params.QueueName == "" {
		params.Log.Warn("User Registration queue name is empty")
		return nil
	}

	registrationQueue, err := params.Registerer.QueueDeclare(
		params.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Log.Error("Cannot create queue", slog.String("queueName", params.QueueName))
		return err
	}
	params.Log.Info("Created userRegistration queue", slog.String("queueName", registrationQueue.Name))

	// Bind Queue to Exchange
	var chosenRoutingKey string
	if params.RoutingKey == "" {
		chosenRoutingKey = defaultAuthRoutingKey
	} else {
		chosenRoutingKey = params.RoutingKey
	}

	if err = params.Registerer.QueueBind(registrationQueue.Name, chosenRoutingKey, AuthExchange, false, nil); err != nil {
		params.Log.Error(
			"Cannot bind registration queue to exchange",
			slog.String("queueName", registrationQueue.Name),
			slog.String("exchangeName", AuthExchange),
		)
		return err
	}

	return nil
}
