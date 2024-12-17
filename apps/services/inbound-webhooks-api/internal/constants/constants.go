package constants

// Event Producer/Consumer
const (
	AuthExchangeName = "authExchange"
	AuthQueueName    = "authQueue"
)

// Event Names for Auth Producer and Consumer
type EventName string

func (n EventName) String() string {
	return string(n)
}

// Event Names
const (
	UserRegistered EventName = "career-cue.auth.userRegistered"
)
