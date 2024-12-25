package user_test

import (
	"libs/backend/domain/user"
	"testing"
)

func TestUserDomain(t *testing.T) {
	result := user.UserDomain("works")
	if result != "UserDomain works" {
		t.Error("Expected UserDomain to append 'works'")
	}
}
