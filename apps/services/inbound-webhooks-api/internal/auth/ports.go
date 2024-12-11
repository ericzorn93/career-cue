package auth

import (
	"apps/services/inbound-webhooks-api/internal/entities"
	"context"
)

// AuthServicePort is the interface that defines
// how the authentication provider interacts with
// the internal system
type AuthServicePort interface {
	// SendUserRegistered will convert the domain type type and send to
	// the message broker
	SendUserRegistered(ctx context.Context, user entities.User) error
}
