package config

import (
	"os"
	"strconv"
)

// Config for the application
type Config struct {
	AMQPUrl string
	RPCPort uint64
}

// NewConfig constructs the config
func NewConfig() (Config, error) {
	port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		return Config{}, err
	}

	config := Config{
		AMQPUrl: os.Getenv("AMQP_CONNECTION_URI"),
		RPCPort: uint64(port),
	}

	return config, nil
}
