package eventing

import (
	boot "libs/backend/boot"
	"log/slog"
)

// Event Producer/Consumer
const (
	AuthExchange = "authExchange"
	AuthDomain   = "auth"
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

func NewAuthEventSetup(registerer boot.AMQPRegisterer, log boot.Logger) AuthEventSetup {
	return AuthEventSetup{
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

func (a *AuthEventSetup) CreateQueue(queueName string) *AuthEventSetup {
	// Create a queue for the service
	if queueName == "" {
		a.log.Warn("Auth queue name is empty")
		return a
	}

	// Register auth queue with the broker
	registrationQueue, err := a.registerer.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
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
func (a AuthEventSetup) Complete() {
	a.log.Info("Completed auth event setup")
}
