package ports

import (
	"context"
	userEntities "libs/backend/domain/user/entities"
)

// AccountRepository is the interface for the account repository
type AccountRepository interface {
	// CreateAccount creates an account in the database
	CreateAccount(context.Context, userEntities.User) error
}
