package config

import (
	"os"
)

// Config for the application
type Config struct {
	AMQPUrl    string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	config := Config{
		AMQPUrl:    os.Getenv("AMQP_CONNECTION_URI"),
		DBHost:     os.Getenv("DATABASE_HOST"),
		DBName:     os.Getenv("DATABASE_NAME"),
		DBUser:     os.Getenv("DATABASE_USER"),
		DBPassword: os.Getenv("DATABASE_PASSWORD"),
		DBPort:     os.Getenv("DATABASE_PORT"),
		DBSSLMode:  os.Getenv("DATABASE_SSL_MODE"),
		DBTimeZone: os.Getenv("DATABASE_TIMEZONE"),
	}

	return config, nil
}
