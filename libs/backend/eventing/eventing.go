package eventing

func Eventing(name string) string {
	result := "Eventing " + name
	return result
}

// Event Names for Auth Producer and Consumer
type EventInfrastructureName string

// String returns the string representation of the event name
func (n EventInfrastructureName) String() string {
	return string(n)
}

// Event Names for Auth Producer and Consumer
type EventName string

// String returns the string representation of the event name
func (n EventName) String() string {
	return string(n)
}
