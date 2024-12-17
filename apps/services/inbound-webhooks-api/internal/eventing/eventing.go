package eventing

import (
	"libs/backend/eventing"
	"libs/boot/pkg/logger"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

// NewAuthQueueParams are params for the auth queue constructor
type NewAuthQueueParams struct {
	Channel *amqp091.Channel
	Log     logger.Logger
}

// RegisterAuthEvents constructs Auth Queue from AMQP Channel
func RegisterAuthEvents(params NewAuthQueueParams) {
	err := params.Channel.ExchangeDeclare(eventing.AuthExchange.String(), "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		os.Exit(1)
	}

	authQueue, err := params.Channel.QueueDeclare(
		eventing.RegistrationQueue.String(),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Log.Error("Cannot create queue")
		os.Exit(1)
	}
	params.Log.Info("Created auth queue")

	if err = params.Channel.QueueBind(authQueue.Name, eventing.RegistrationQueue.String(), eventing.AuthExchange.String(), false, nil); err != nil {
		params.Log.Error("Cannot bind auth queue")
		os.Exit(1)
		return
	}
}
