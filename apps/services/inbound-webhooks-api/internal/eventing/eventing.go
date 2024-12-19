package eventing

import (
	"libs/backend/eventing"
	"libs/boot/pkg/amqp"
	"libs/boot/pkg/logger"
	"os"
)

// NewAuthQueueParams are params for the auth queue constructor
type NewAuthQueueParams struct {
	Registerer amqp.Registerer
	Log        logger.Logger
}

// RegisterAuthEvents constructs Auth Queue from AMQP Channel
func RegisterAuthEvents(params NewAuthQueueParams) {
	err := params.Registerer.ExchangeDeclare(eventing.AuthExchange.String(), "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		os.Exit(1)
	}

	authQueue, err := params.Registerer.QueueDeclare(
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

	if err = params.Registerer.QueueBind(authQueue.Name, eventing.RegistrationQueue.String(), eventing.AuthExchange.String(), false, nil); err != nil {
		params.Log.Error("Cannot bind auth queue")
		os.Exit(1)
		return
	}
}
