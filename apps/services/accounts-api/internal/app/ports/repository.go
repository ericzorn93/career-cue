package ports

import (
	"context"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
	"time"
)

// AccountRepository is the interface for the account repository
type AccountRepository interface {
	// CreateAccount creates an account in the database
	CreateAccount(context.Context, userEntities.User) error

	// GetAccountByCommonID gets an account from the database by commonID
	GetAccountByCommonID(context.Context, userValueObjects.CommonID) (userEntities.User, error)

	// GetAccountByEmailAddress gets an account from the database by email address
	GetAccountByEmailAddress(context.Context, userValueObjects.EmailAddress) (userEntities.User, error)

	// SoftDeleteAccountByCommonID will soft delete the user
	SoftDeleteAccountByCommonID(context.Context, userValueObjects.CommonID) (time.Time, error)

	// HardDeleteAccountByCommonID will hard delete the user
	HardDeleteAccountByCommonID(context.Context, userValueObjects.CommonID) (time.Time, error)
}
