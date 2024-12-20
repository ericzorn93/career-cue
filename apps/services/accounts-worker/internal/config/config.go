package config

import (
	"fmt"
	"os"
)

const (
	// serviceName is the name of the microservice
	serviceName = "accounts-worker"

	// userRegistrationQueueName is the name of the queue for user registration
	userRegistrationQueueName = "user-registration"
)

// Config for the application
type Config struct {
	ServiceName               string
	AMQPUrl                   string
	UserRegistrationQueueName string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		ServiceName:               serviceName,
		AMQPUrl:                   os.Getenv("AMQP_CONNECTION_URI"),
		UserRegistrationQueueName: fmt.Sprintf("%s-%s", serviceName, userRegistrationQueueName),
	}

	return config, nil
}
