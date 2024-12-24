package ports

import "apps/services/inbound-webhooks-api/internal/domain/entities"

// AuthService will handle auth webhook
// interactions
type AuthService interface {
	RegisterUser(user entities.User) error
}
