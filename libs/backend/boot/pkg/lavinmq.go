package boot

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// LavinMQOptions configuration to start
// the LavinMQ connections to queues and exchanges
type LavinMQOptions struct {
	ConnectionURI        string
	OnConnectionCallback func(conn *amqp.Connection) error
}
