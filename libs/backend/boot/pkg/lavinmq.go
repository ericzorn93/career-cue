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
	lavinMQModuleName = "lavinMQ"
)

// LavinMQOutput will output each individual property
type LavinMQOutput struct {
	fx.Out

	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewLavinMQModule() fx.Option {
	return fx.Module(
		lavinMQModuleName,
		fx.Provide(func(lc fx.Lifecycle, bsParams BootServiceParams, log *slog.Logger) (LavinMQOutput, error) {
			// Validate the conneciton URI to LavinMQ
			if bsParams.LavinMQOptions.ConnectionURI == "" || !strings.HasPrefix(bsParams.LavinMQOptions.ConnectionURI, "amqp://") {
				log.Error("Cannot connect to LavinMQ with empty connection string")
				return LavinMQOutput{}, errors.New("cannot connect with invalid LavinMQ String")
			}

			conn, err := amqp.DialConfig(bsParams.LavinMQOptions.ConnectionURI, amqp.Config{
				Heartbeat: time.Second * 10,
			})
			if err != nil {
				log.Error("Cannot connect to LavinMQ", slog.Any("error", err))
				return LavinMQOutput{}, errors.New("LavinMQ connection failed")
			}

			// Start/Stop the connection on close
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := bsParams.LavinMQOptions.OnConnectionCallback(); err != nil {
						log.Error("LavinMQ connection callback failed", slog.Any("error", err))
						return err
					}

					return nil
				},
				OnStop: func(_ context.Context) error {
					log.Info("Closing the LavinMQ connection")
					return conn.Close()
				},
			})

			// Create Channel
			ch, err := conn.Channel()
			if err != nil {
				log.Error("Cannot create channel")
				return LavinMQOutput{}, err
			}

			return LavinMQOutput{Conn: conn, Channel: ch}, nil
		}),
		fx.Invoke(func(_ *amqp.Connection, _ *amqp.Channel, log *slog.Logger) {
			log.Info("LavinMQ connection and channel are ready for use")
		}),
	)
}

// LavinMQOptions configuration to start
// the LavinMQ connections to queues and exchanges
type LavinMQOptions struct {
	ConnectionURI        string
	OnConnectionCallback func() error
}
