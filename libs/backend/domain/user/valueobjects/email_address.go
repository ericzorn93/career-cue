package valueobjects

// EmailAddress is a value object for an email address
type EmailAddress struct {
	value string
}

// NewEmailAddress creates a new EmailAddress
func NewEmailAddress(value string) EmailAddress {
	return EmailAddress{value: value}
}

// Value returns the value of the EmailAddress
func (e EmailAddress) Value() string {
	return e.value
}

// Equals checks if two EmailAddress objects are equal
func (e EmailAddress) Equals(other EmailAddress) bool {
	return e.value == other.value
}

// String returns the string representation of the EmailAddress
func (e EmailAddress) String() string {
	return e.value
}

// IsEmpty checks if the EmailAddress is empty
func (e EmailAddress) IsEmpty() bool {
	return e.value == ""
}
