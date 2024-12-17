package amqp

import (
	"errors"
	"libs/boot/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// CallBackParams returns the logger and callback and channel to the caller
type CallBackParams struct {
	Logger     logger.Logger
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// Options configuration to start
// the LavinMQ connections to queues and exchanges
type Options struct {
	ConnectionURI        string
	OnConnectionCallback func(CallBackParams) error
}

// EstablishAMQPConnection will create a connection to the AMQP broker
func EstablishAMQPConnection(log logger.Logger, opts Options) (*amqp.Connection, *amqp.Channel, error) {
	// Validate the conneciton URI to AMQP
	if opts.ConnectionURI == "" || !strings.HasPrefix(opts.ConnectionURI, "amqp://") {
		log.Error("Cannot connect to AMQP with empty connection string")
		return nil, nil, errors.New("cannot connect with invalid AMQP String")
	}

	conn, err := amqp.DialConfig(opts.ConnectionURI, amqp.Config{
		Heartbeat: time.Second * 10,
	})
	if err != nil {
		log.Error("Cannot connect to AMQP", slog.Any("error", err))
		return nil, nil, errors.New("AMQP connection failed")
	}

	// Handler server close
	go func() {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
		<-exitCh

		log.Info("Closing the AMQP connection")
		if err := conn.Close(); err != nil {
			log.Error("Trouble closing the AMQP Connection")
		}
	}()

	// Create Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error("Cannot create AMQP channel")
		return nil, nil, err
	}

	// Start/Stop the connection on close
	if err := opts.OnConnectionCallback(CallBackParams{
		Logger:     log,
		Connection: conn,
		Channel:    ch,
	}); err != nil {
		log.Error("AMQP connection callback failed", slog.Any("error", err))
		return nil, nil, err
	}

	return conn, ch, err
}
