package boot

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// LavinMQOptions configuration to start
// the LavinMQ connections to queues and exchanges
type LavinMQOptions struct {
	ConnectionURI        string
	OnConnectionCallback func() error
}

// LavinMQ is the parent struct that wraps
// the lavinMQ connection pointer
type LavinMQ struct {
	Connection *amqp.Connection
}
