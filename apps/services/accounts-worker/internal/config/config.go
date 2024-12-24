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
	AMQPUri                   string
	UserRegistrationQueueName string
	AccountsAPIUri            string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		ServiceName:               serviceName,
		AMQPUri:                   os.Getenv("AMQP_CONNECTION_URI"),
		UserRegistrationQueueName: fmt.Sprintf("%s-%s", serviceName, userRegistrationQueueName),
		AccountsAPIUri:            os.Getenv("ACCOUNTS_API_URI"),
	}

	return config, nil
}
