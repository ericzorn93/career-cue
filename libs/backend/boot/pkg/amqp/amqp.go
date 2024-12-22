package amqp

import (
	"context"
	"errors"
	"libs/boot/pkg/logger"
	"log/slog"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// CallBackParams returns the logger and callback and channel to the caller
type CallBackParams struct {
	Logger     logger.Logger
	Controller Controller
}

type HandlerParams struct {
	Logger         logger.Logger
	AMQPController Controller
}

type Handler func(HandlerParams) error

// Options configuration to start
// the LavinMQ connections to queues and exchanges
type Options struct {
	ConnectionURI        string
	OnConnectionCallback func(CallBackParams) error
	Handlers             []Handler
}

func (o Options) IsZero() bool {
	return o.ConnectionURI == "" && o.OnConnectionCallback == nil
}

// EstablishAMQPConnection will create a connection to the AMQP broker
func EstablishAMQPConnection(log logger.Logger, opts Options) (Controller, error) {
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
		log.Error("Cannot connect to AMQP with empty connection string")
		return Controller{}, errors.New("cannot connect with invalid AMQP String")
	}

	// Connect to AMQP broker and poll every 10 seconds god health
	conn, err := amqp.DialConfig(opts.ConnectionURI, amqp.Config{
		Heartbeat: time.Second * 10,
	})
	if err != nil {
		log.Error("Cannot connect to AMQP", slog.Any("error", err))
		return Controller{}, errors.New("AMQP connection failed")
	}

	// // Handler server close
	// go func() {
	// 	exitCh := make(chan os.Signal, 1)
	// 	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
	// 	<-exitCh

	// 	log.Info("Closing the AMQP connection")
	// 	if err := conn.Close(); err != nil {
	// 		log.Error("Trouble closing the AMQP Connection")
	// 	}
	// }()

	// Create Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error("Cannot create AMQP channel")
		return Controller{}, err
	}

	// AMQP controller wrapper
	controller := NewController(ch)

	// Start/Stop the connection on close
	callbackParams := CallBackParams{
		Logger:     log,
		Controller: controller,
	}
	if err := opts.OnConnectionCallback(callbackParams); err != nil {
		log.Error("AMQP connection callback failed", slog.Any("error", err))
		return controller, err
	}

	// Register All Handlers
	var wg sync.WaitGroup
	if len(opts.Handlers) == 0 {
		log.Warn("No AMQP handlers present")
	} else {
		handlerParams := HandlerParams{
			Logger:         log,
			AMQPController: controller,
		}
		wg.Add(len(opts.Handlers))
		for _, handler := range opts.Handlers {
			handler := handler
			go func() {
				forever := make(chan struct{})
				if err := handler(handlerParams); err != nil {
					log.Error("Cannot register AMQP handler", slog.Any("error", err))
					close(forever)
				}
				wg.Done()
				<-forever
			}()
		}
		wg.Wait()
	}

	return controller, nil
}

// Publisher defines the AMQP publish methods
type Publisher interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	PublishWithContext(_ context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// Consumer defines the AMQP consume methods
type Consumer interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

// Registerer defines the AMQP register methods for queues and exchanges
type Registerer interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

// Controller returns an interface for publishing, consuming and registering
type Controller struct {
	channel    *amqp.Channel
	Publisher  Publisher
	Consumer   Consumer
	Registerer Registerer
}

// NewController constructs the returns object for controlling AMQP
func NewController(channel *amqp.Channel) Controller {
	return Controller{
		Publisher:  channel,
		Consumer:   channel,
		Registerer: channel,
	}
}

// IsConnected will let the caller know if the controller has established an AMQP broker connection
func (c Controller) IsConnected() bool {
	return c.channel == nil || c.channel.IsClosed()
}

// Close will close the AMQP connection
func (c Controller) Close() error {
	return c.channel.Close()
}
