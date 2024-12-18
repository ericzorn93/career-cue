package eventing

import (
	"testing"
)

func TestEventing(t *testing.T) {
	result := Eventing("works")
	if result != "Eventing works" {
		t.Error("Expected Eventing to append 'works'")
	}
}
