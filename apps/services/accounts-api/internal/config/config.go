package config

import (
	"os"
)

// Config for the application
type Config struct {
	AMQPUrl string
	DBDsn   string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		AMQPUrl: os.Getenv("AMQP_CONNECTION_URI"),
		DBDsn:   os.Getenv("DB_CONNECTION_DSN"),
	}

	return config, nil
}
