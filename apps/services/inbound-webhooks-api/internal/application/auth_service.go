package application

import "apps/services/inbound-webhooks-api/internal/domain"

// AuthService will handle auth webhook
// interactions
type AuthService interface {
	RegisterUser(user domain.User) error
}
