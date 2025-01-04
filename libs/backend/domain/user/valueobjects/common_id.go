package valueobjects

import (
	"github.com/google/uuid"
)

// CommonID is a value object for a common ID
type CommonID struct {
	value uuid.UUID
}

// NewCommonID creates a new CommonID
func NewCommonID() CommonID {
	return CommonID{value: uuid.New()}
}

// NewCommonID creates a new CommonID
func NewCommonIDFromString(value string) CommonID {
	commonID, err := uuid.Parse(value)
	if err != nil {
		return CommonID{}
	}

	return CommonID{value: commonID}
}

// NewCommonIDFromUUID creates a new CommonID from uuid.UUID instance
func NewCommonIDFromUUID(id uuid.UUID) CommonID {
	return CommonID{value: id}
}

// Value returns the value of the CommonID
func (c CommonID) Value() uuid.UUID {
	return c.value
}

// Equals checks if two CommonID objects are equal
func (c CommonID) Equals(other CommonID) bool {
	return c.value == other.value
}

// String returns the string representation of the CommonID
func (c CommonID) String() string {
	return c.value.String()
}

// IsEmtpy checks if the CommonID is empty
func (c CommonID) IsEmpty() bool {
	return c.value == uuid.Nil
}
