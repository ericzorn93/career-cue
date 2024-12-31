package ports

import (
	"context"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
)

// AccountRepository is the interface for the account repository
type AccountRepository interface {
	// CreateAccount creates an account in the database
	CreateAccount(context.Context, userEntities.User) error

	// GetAccount gets an account from the database by commonID
	GetAccount(context.Context, userValueObjects.CommonID) (userEntities.User, error)
}
