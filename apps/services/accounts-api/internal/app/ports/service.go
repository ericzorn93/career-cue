package ports

import (
	"context"
	userEntities "libs/backend/domain/user/entities"
)

// RegistrationService is the interface for the registration service
type RegistrationService interface {
	// RegisterUser registers a user in the system and the database
	RegisterUser(ctx context.Context, user userEntities.User) error
}
