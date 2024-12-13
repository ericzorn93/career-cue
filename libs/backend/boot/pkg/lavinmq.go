package boot

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

const (
	moduleName = "lavinMQ"
)

func NewLavinMQModule() fx.Option {
	return fx.Module(
		moduleName,
		fx.Provide(func(lc fx.Lifecycle, bsParams BootServiceParams, log *slog.Logger) (*amqp.Connection, error) {
			// Validate the conneciton URI to LavinMQ
			if bsParams.LavinMQOptions.ConnectionURI == "" || !strings.HasPrefix(bsParams.LavinMQOptions.ConnectionURI, "amqp://") {
				log.Error("Cannot connect to LavinMQ with empty connection string")
				return nil, errors.New("cannot connect with invalid LavinMQ String")
			}

			conn, err := amqp.DialConfig(bsParams.LavinMQOptions.ConnectionURI, amqp.Config{
				Heartbeat: time.Second * 10,
			})
			if err != nil {
				log.Error("Cannot connect to LavinMQ", slog.Any("error", err))
				return conn, errors.New("LavinMQ connection failed")
			}

			// Stop the connection on close
			lc.Append(fx.Hook{
				OnStop: func(_ context.Context) error {
					log.Info("Closing the LavinMQ connection")
					return conn.Close()
				},
			})

			// Start LavinMQ Connection
			if err := bsParams.LavinMQOptions.OnConnectionCallback(); err != nil {
				log.Error("LavinMQ connection callback failed", slog.Any("error", err))
				return nil, errors.New("LavinMQ callback failed")
			}

			return conn, nil
		}),
		fx.Invoke(func(conn *amqp.Connection, log *slog.Logger) {
			log.Info("LavinMQ connection is ready for use")
			// Test the connection
			ch, err := conn.Channel()
			if err != nil {
				log.Error("Failed to open a channel", slog.Any("error", err))
				return
			}
			defer ch.Close()

			log.Info("Successfully opened a channel")
		}),
	)
}

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
