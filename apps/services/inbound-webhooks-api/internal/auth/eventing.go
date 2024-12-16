package auth

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

// NewAuthQueueParams are params for the auth queue constructor
type NewAuthQueueParams struct {
	fx.In

	Channel *amqp.Channel
	Log     *slog.Logger
}

// NewAuthQueue constructs Auth Queue from AMQP Channel
func NewAuthQueue(params NewAuthQueueParams) (amqp.Queue, error) {
	const authQueueName = "authQueue"

	q, err := params.Channel.QueueDeclare(
		authQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Log.Error("Cannot create queue")
		return amqp.Queue{}, err
	}

	params.Log.Info("Created auth queue")
	return q, nil
}
