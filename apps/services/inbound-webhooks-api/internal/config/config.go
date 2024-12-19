package config

import (
	"os"
)

// Config for the application
type Config struct {
	AMQPUrl               string
	RegistrationQueueName string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {

	config := Config{
		AMQPUrl:               os.Getenv("AMQP_CONNECTION_URI"),
		RegistrationQueueName: os.Getenv("AUTH_REGISTRATION_QUEUE_NAME"),
	}

	return config, nil
}
