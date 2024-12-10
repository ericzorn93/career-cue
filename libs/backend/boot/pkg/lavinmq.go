package boot

// LavinMQOptions configuration to start
// the LavinMQ Connections
type LavinMQOptions struct {
	ConnectionURI string
	Queues        []LavinMQQueueOptions
	Exchanges     []LavinMQExchangeOptions
}

// LavinMQExchangeOptions sets up options for exhange
type LavinMQExchangeOptions struct {
	Name         string
	ExchangeType string
}

// LavinMQQueueOptions sets up options for queue
type LavinMQQueueOptions struct {
	Name              string
	QueueType         string
	IsDeadLetter      bool
	ExchangeBoundName string
}
