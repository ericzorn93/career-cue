package auth

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

// NewChannelParams for registering the new channel
type NewChannelParams struct {
	fx.In
	LC     fx.Lifecycle
	Logger *slog.Logger
	Conn   *amqp.Connection
}

// NewChannel registers the connection and channel with Fx
func NewChannel(params NewChannelParams) (*amqp.Channel, error) {
	channel, err := params.Conn.Channel()
	if err != nil {
		return &amqp.Channel{}, err
	}

	params.LC.Append(fx.StopHook(func() error {
		params.Logger.Error("Closing the LavinMQ channel")
		return channel.Close()
	}))

	return channel, nil
}

// NewAuthQueueParams will add new auth
// queue dependencies to the constructor
type NewAuthQueueParams struct {
	fx.In
	LC      fx.Lifecycle
	Logger  *slog.Logger
	Channel *amqp.Channel
}

// NewAuthQueue sets up auth queue
func NewAuthQueue(params NewAuthQueueParams) (amqp.Queue, error) {
	const queueName = "authQueue"
	queue, err := params.Channel.QueueDeclare(
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
