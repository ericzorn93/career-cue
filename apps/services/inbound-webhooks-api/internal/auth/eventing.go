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

// NewAuthQueue constructs Auth Queue from AMQP Channel
func NewAuthQueue(params NewAuthQueueParams) (func(message []byte) error, error) {
	const authExchangeName = "authExchange"
	const authQueueName = "authQueue"

	err := params.Channel.ExchangeDeclare(authExchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		params.Log.Error("Cannot create exchange")
		return nil, err
	}

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
		return nil, err
	}
	params.Log.Info("Created auth queue")

	if err = params.Channel.QueueBind(q.Name, "", authExchangeName, false, nil); err != nil {
		return nil, err
	}

	return func(message []byte) error {
		return params.Channel.Publish(authExchangeName, "", false, false, amqp.Publishing{Body: message})
	}, nil
}
