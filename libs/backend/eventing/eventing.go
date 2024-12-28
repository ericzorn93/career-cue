package eventing

import "fmt"

func Eventing(name string) string {
	result := "Eventing " + name
	return result
}

// Event Names for Auth Producer and Consumer
const (
	CareerCueEventPrefix = "career-cue"
)

type EventName string

// String returns the string representation of the event name
func (n EventName) String() string {
	return string(n)
}

// GetRoutingKeyPrefix concats the app name and domain name
func GetRoutingKeyPrefix(domainName string) string {
	return fmt.Sprintf("%s.%s", CareerCueEventPrefix, domainName)
}

// GetEventName will construct the event name from the app name, domain and specific event name
func GetEventName(domain, eventName string) string {
	return fmt.Sprintf("%s.%s.%s", CareerCueEventPrefix, domain, eventName)
}

// EventBuilder is an interface for building events and event infrastructure
type EventBuilder interface {
	CreateExchange() EventBuilder
	CreateDeadletter() EventBuilder
	CreateQueue(queueName string) EventBuilder
	BindQueues(routingKeys []string) EventBuilder
	Complete()
}
