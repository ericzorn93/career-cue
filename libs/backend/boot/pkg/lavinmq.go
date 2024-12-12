package boot

// LavinMQOptions configuration to start
// the LavinMQ connections to queues and exchanges
type LavinMQOptions struct {
	ConnectionURI        string
	OnConnectionCallback func() error
}
