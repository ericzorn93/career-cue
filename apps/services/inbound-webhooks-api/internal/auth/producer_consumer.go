package auth

import (
	boot "libs/boot/pkg"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

// EstablishQueuesParams will add new auth
// queue dependencies to the constructor
type EstablishQueuesParams struct {
	fx.In
	LC          fx.Lifecycle
	Logger      *slog.Logger
	LavinMQConn boot.LavinMQ
}

// EstablishQueues sets up auth queue
func EstablishQueues(params EstablishQueuesParams) (amqp.Queue, error) {
	// Create channel and closer function
	channel, err := params.LavinMQConn.Connection.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}

	params.LC.Append(fx.StopHook(func() error {
		params.Logger.Error("Closing the LavinMQ channel")
		return channel.Close()
	}))

	// Establish Auth Queue
	const queueName = "authQueue"
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Logger.Error("Could not create the auth queue")
		return amqp.Queue{}, err
	}

	return queue, nil
}
