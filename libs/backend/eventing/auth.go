package eventing

import (
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"
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
	authRoutingKey = GetRoutingKeyPrefix(AuthDomain) + ".*"
)

// RegisterAuthParams are params for the auth queue constructor
type RegisterAuthParams struct {
	Registerer amqp.Registerer
	Log        logger.Logger
	QueueName  string
}

// RegisterAuth constructs Auth Queue from AMQP Channel
func RegisterAuth(params RegisterAuthParams) error {
	// Initialize Auth Exchange - topic
	err := params.Registerer.ExchangeDeclare(AuthExchange, "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		return err
	}

	// Create a queue for the service
	if params.QueueName == "" {
		params.Log.Warn("Auth queue name is empty")
		return nil
	}

	authQueue, err := params.Registerer.QueueDeclare(
		params.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Log.Error("Cannot create auth queue", slog.String("queueName", params.QueueName))
		return err
	}
	params.Log.Info("Created auth queue", slog.String("queueName", authQueue.Name))

	// Bind Queue to Exchange
	if err = params.Registerer.QueueBind(authQueue.Name, authRoutingKey, AuthExchange, false, nil); err != nil {
		params.Log.Error(
			"Cannot bind queue to exchange",
			slog.String("queueName", authQueue.Name),
			slog.String("exchangeName", AuthExchange),
		)
		return err
	}

	return nil
}
