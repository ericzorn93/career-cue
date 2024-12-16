package auth

import (
	boot "libs/boot/pkg"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

// NewAuthQueueParams are params for the auth queue constructor
type NewAuthQueueParams struct {
	fx.In

	Channel *amqp.Channel
	Log     boot.Logger
}

// NewAuthPublisher constructs Auth Queue from AMQP Channel
func NewAuthPublisher(params NewAuthQueueParams) (func([]byte) error, error) {
	const authExchangeName = "authExchange"
	const authQueueName = "authQueue"

	err := params.Channel.ExchangeDeclare(authExchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		return nil, err
	}

	authQueue, err := params.Channel.QueueDeclare(
		authQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		params.Log.Error("Cannot create queue")
		return nil, err
	}
	params.Log.Info("Created auth queue")

	if err = params.Channel.QueueBind(authQueue.Name, "", authExchangeName, false, nil); err != nil {
		return nil, err
	}

	return func(msg []byte) error {
		return params.Channel.Publish(authExchangeName, authQueue.Name, false, false, amqp.Publishing{Body: msg})
	}, nil
}
