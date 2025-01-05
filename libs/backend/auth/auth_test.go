package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	result := Auth("works")
	if result != "Auth works" {
		t.Error("Expected Auth to append 'works'")
	}
}
