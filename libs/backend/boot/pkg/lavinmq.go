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

type AmqpPublisher interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type AmqpConsumer interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type AmqpRegisterer interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

type Publisher struct {
	Channel *amqp.Channel
}

func (p Publisher) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return p.Channel.Publish(exchange, key, mandatory, immediate, msg)
}
func (p Publisher) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return p.Channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

type Consumer struct {
	Channel *amqp.Channel
}

func (c Consumer) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
func (c Consumer) ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.Channel.ConsumeWithContext(ctx, queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

type Registerer struct {
	Channel *amqp.Channel
}

func (r Registerer) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return r.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}
func (r Registerer) ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error {
	return r.Channel.ExchangeBind(destination, key, source, noWait, args)
}
func (r Registerer) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}
func (r Registerer) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return r.Channel.QueueBind(name, key, exchange, noWait, args)
}

// LavinMQOutput will output each individual property
type LavinMQOutput struct {
	fx.Out

	Publisher  AmqpPublisher
	Consumer   AmqpConsumer
	Registerer AmqpRegisterer
}

func NewLavinMQModule() fx.Option {
	return fx.Module(
		lavinMQModuleName,
		fx.Provide(func(lc fx.Lifecycle, bsParams BootServiceParams, log Logger) (LavinMQOutput, error) {
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

			return LavinMQOutput{
				Publisher:  Publisher{ch},
				Consumer:   Consumer{ch},
				Registerer: Registerer{ch},
			}, nil
		}),
		fx.Invoke(func(_ AmqpPublisher, _ AmqpConsumer, _ AmqpRegisterer, log Logger) {
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
