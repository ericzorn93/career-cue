package ports

import (
	"context"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
)

// RegistrationService is the interface for the registration service
type RegistrationService interface {
	// RegisterUser registers a user in the system and the database
	RegisterUser(ctx context.Context, user userEntities.User) error

	// GetUser gets a user from the system
	GetUser(ctx context.Context, commonID userValueObjects.CommonID) (userEntities.User, error)
}
