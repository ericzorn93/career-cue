package eventing

import (
	boot "libs/backend/boot"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

// Event Producer/Consumer
const (
	AuthDomain               = "auth"
	AuthExchange             = "authExchange"
	AuthDeadletterExchange   = "authDeadletterExchange"
	AuthDeadletterQueue      = "authDeadletterQueue"
	AuthDeadletterRoutingKey = "authDlx"
)

// Event Names
var (
	EventNameUserRegistered EventName = EventName(GetEventName(AuthDomain, "userRegistered"))
)

// GetUserRegisteredRoutingKey returns the routing key for user registered event
func GetUserRegisteredRoutingKey() string {
	return EventNameUserRegistered.String()
}

// RegisterAuthParams are params for the auth queue constructor
type RegisterAuthParams struct {
	Registerer boot.AMQPRegisterer
	Log        boot.Logger
	QueueName  string
	RoutingKey string
}

// AuthRegisterer is a struct for auth queues and exchanges constructor
type AuthEventSetup struct {
	registerer boot.AMQPRegisterer
	log        boot.Logger
	queueNames []string
}

func NewAuthEventSetup(registerer boot.AMQPRegisterer, log boot.Logger) *AuthEventSetup {
	return &AuthEventSetup{
		registerer: registerer,
		log:        log,
		queueNames: make([]string, 0),
	}
}

// CreateExchange creates an exchange for the auth event setup
func (a *AuthEventSetup) CreateExchange() *AuthEventSetup {
	// Initialize Auth Exchange - topic
	err := a.registerer.ExchangeDeclare(AuthExchange, "topic", true, false, false, false, nil)
	if err != nil {
		a.log.Error("Cannot create exchange")
		return a
	}

	a.log.Info("Created auth exchange", slog.String("exchangeName", AuthExchange))
	return a
}

// CreateDeadletter creates a dead letter exchange and queue for the auth event setup
func (a *AuthEventSetup) CreateDeadletter() *AuthEventSetup {
	// Initialize Dead Letter Exchange
	err := a.registerer.ExchangeDeclare(AuthDeadletterExchange, "direct", true, false, false, false, nil)
	if err != nil {
		a.log.Error("Cannot create dead letter exchange")
		return a
	}

	// Initialize Dead Letter Queue
	authDeadletterQueue, err := a.registerer.QueueDeclare(
		AuthDeadletterQueue, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		a.log.Error("Failed to declare the auth dead letter queue:")
		return a
	}

	// Bind Auth Dead Letter Queue to Auth Dead Letter Exchange
	err = a.registerer.QueueBind(
		authDeadletterQueue.Name, // queue name
		AuthDeadletterRoutingKey, // routing key
		AuthDeadletterExchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		a.log.Error("Cannot bind auth deadletter queue to exchange")
		return a
	}

	a.log.Info("Created auth dead letter exchange and queue", slog.String("exchangeName", AuthDeadletterExchange), slog.String("queueName", authDeadletterQueue.Name))

	return a
}

func (a *AuthEventSetup) CreateQueue(queueName string) *AuthEventSetup {
	// Create a queue for the service
	if queueName == "" {
		a.log.Warn("Auth queue name is empty")
		return a
	}

	// Register auth queue with the broker
	registrationQueue, err := a.registerer.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp091.Table{
			"x-dead-letter-exchange":    AuthDeadletterExchange,
			"x-dead-letter-routing-key": AuthDeadletterRoutingKey,
		},
	)
	if err != nil {
		a.log.Error("Cannot create queue", slog.String("queueName", queueName))
		return a
	}

	a.log.Info("Created auth queue", slog.String("queueName", registrationQueue.Name))

	// Add queue name to the list of queues for auth
	a.queueNames = append(a.queueNames, registrationQueue.Name)

	return a
}

// BindQueues binds the queues to the exchange for each routing key
func (a *AuthEventSetup) BindQueues(routingKeys []string) *AuthEventSetup {
	a.log.Info("Binding queues to exchange", slog.String("exchangeName", AuthExchange), slog.Any("queueNames", a.queueNames))

	// Bind Queue to Exchange
	for _, queueName := range a.queueNames {
		for _, routingKey := range routingKeys {
			if err := a.registerer.QueueBind(queueName, routingKey, AuthExchange, false, nil); err != nil {
				a.log.Error(
					"Cannot bind registration queue to exchange",
					slog.String("queueName", queueName),
					slog.String("exchangeName", AuthExchange),
				)
				return a
			}
		}
	}

	return a
}

// Complete completes the auth event setup
func (a *AuthEventSetup) Complete() {
	a.log.Info("Completed auth event setup")
}
