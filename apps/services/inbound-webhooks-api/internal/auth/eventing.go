package auth

import (
	boot "libs/boot/pkg"
	"os"

	"go.uber.org/fx"
)

// NewAuthQueueParams are params for the auth queue constructor
type NewAuthQueueParams struct {
	fx.In

	Registerer boot.AmqpRegisterer
	Log        boot.Logger
}

// RegisterAuthEvents constructs Auth Queue from AMQP Channel
func RegisterAuthEvents(params NewAuthQueueParams) {
	err := params.Registerer.ExchangeDeclare(AuthExchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		os.Exit(1)
	}

	authQueue, err := params.Registerer.QueueDeclare(
		AuthQueueName,
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

	if err = params.Registerer.QueueBind(authQueue.Name, "", AuthExchangeName, false, nil); err != nil {
		params.Log.Error("Cannot bind auth queue")
		os.Exit(1)
		return
	}
}
