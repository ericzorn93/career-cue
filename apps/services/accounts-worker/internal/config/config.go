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

	Auth0Domain       string
	Auth0ClientID     string
	Auth0ClientSecret string
	Auth0Audience     string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		ServiceName:               serviceName,
		AMQPUri:                   os.Getenv("AMQP_CONNECTION_URI"),
		UserRegistrationQueueName: fmt.Sprintf("%s-%s", serviceName, userRegistrationQueueName),
		AccountsAPIUri:            os.Getenv("ACCOUNTS_API_URI"),

		Auth0Domain:       os.Getenv("AUTH0_DOMAIN"),
		Auth0ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		Auth0Audience:     os.Getenv("AUTH0_AUDIENCE"),
	}

	return config, nil
}
