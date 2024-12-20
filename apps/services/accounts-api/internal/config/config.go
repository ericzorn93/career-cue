package config

import (
	"os"
)

// Config for the application
type Config struct {
	AMQPUrl string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		AMQPUrl: os.Getenv("AMQP_CONNECTION_URI"),
	}

	return config, nil
}
