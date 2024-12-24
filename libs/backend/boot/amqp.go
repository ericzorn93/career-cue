package boot

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPCallBackParams returns the logger and callback and channel to the caller
type AMQPCallBackParams struct {
	Logger     Logger
	Controller AMQPController
}

// AMQPHandlerParams returns the logger and AMQP controller to the handler
type AMQPHandlerParams struct {
	Logger         Logger
	AMQPController AMQPController
}

// AMQPHandler is a type of callback used specifically for starting the AMQP handlers
type AMQPHandler func(AMQPHandlerParams) error

// AMQPOptions configuration to start
// the LavinMQ connections to queues and exchanges
type AMQPOptions struct {
	ConnectionURI        string
	OnConnectionCallback func(AMQPCallBackParams) error
	Handlers             []AMQPHandler
}

// IsZero will let the caller know if the AMQPOptions is empty
func (o AMQPOptions) IsZero() bool {
	return o.ConnectionURI == "" && o.OnConnectionCallback == nil
}

// EstablishAMQPConnection will create a connection to the AMQP broker
func (s *BootService) EstablishAMQPConnection(opts AMQPOptions) error {
	// Validate the conneciton URI to AMQP
	connectionPrefixes := [2]string{"amqp://", "amqps://"}
	hasPrefix := false
	for _, prefix := range connectionPrefixes {
		hasPrefix = strings.HasPrefix(opts.ConnectionURI, prefix)
		if hasPrefix {
			break
		}
	}

	if opts.ConnectionURI == "" || !hasPrefix {
		s.logger.Error("Cannot connect to AMQP with empty connection string")
		return errors.New("cannot connect with invalid AMQP String")
	}

	// Connect to AMQP broker and poll every 10 seconds god health
	conn, err := amqp.DialConfig(opts.ConnectionURI, amqp.Config{
		Heartbeat: time.Second * 10,
	})
	if err != nil {
		s.logger.Error("Cannot connect to AMQP", slog.Any("error", err))
		return errors.New("AMQP connection failed")
	}

	// Create Channel
	ch, err := conn.Channel()
	if err != nil {
		s.logger.Error("Cannot create AMQP channel")
		return err
	}

	// AMQP controller wrapper
	controller := NewController(s.logger, conn, ch)
	s.amqpController = controller

	// Start/Stop the connection on close
	callbackParams := AMQPCallBackParams{
		Logger:     s.logger,
		Controller: controller,
	}
	if err := opts.OnConnectionCallback(callbackParams); err != nil {
		s.logger.Error("AMQP connection callback failed", slog.Any("error", err))
		return err
	}

	// Register All Handlers
	if len(opts.Handlers) == 0 {
		s.logger.Warn("No AMQP handlers present")
		return nil
	}

	s.logger.Info("Registering AMQP handlers", slog.Int("count", len(opts.Handlers)))

	// Register all the handlers and run indefinitely
	forever := make(chan struct{})
	for _, handler := range opts.Handlers {
		handlerParams := AMQPHandlerParams{
			Logger:         s.logger,
			AMQPController: controller,
		}
		go handler(handlerParams)
	}
	<-forever
	return nil
}

// AMQPPublisher defines the AMQP publish methods
type AMQPPublisher interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	PublishWithContext(_ context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// AMQPConsumer defines the AMQP consume methods
type AMQPConsumer interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

// AMQPRegisterer defines the AMQP register methods for queues and exchanges
type AMQPRegisterer interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

// AMQPController returns an interface for publishing, consuming and registering
type AMQPController struct {
	logger     Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	Publisher  AMQPPublisher
	Consumer   AMQPConsumer
	Registerer AMQPRegisterer
}

// NewController constructs the returns object for controlling AMQP
func NewController(logger Logger, connection *amqp.Connection, channel *amqp.Channel) AMQPController {
	return AMQPController{
		logger:     logger,
		connection: connection,
		channel:    channel,
		Publisher:  channel,
		Consumer:   channel,
		Registerer: channel,
	}
}

// IsConnected will let the caller know if the controller has established an AMQP broker connection
func (c AMQPController) IsConnected() bool {
	return c.connection != nil && !c.connection.IsClosed()
}

// Close will close the AMQP connection
func (c AMQPController) Close() error {
	if !c.IsConnected() {
		c.logger.Warn("AMQP connection is already closed")
		return nil
	}

	c.logger.Info("Closing the AMQP connection")
	return c.channel.Close()
}
