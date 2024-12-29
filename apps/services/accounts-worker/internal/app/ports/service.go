package ports

import (
	"context"
	"errors"
	userEntities "libs/backend/domain/user/entities"
)

var (
	// ErrUserAlreadyExists is returned when a user already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserNotCreated is returned when a user is not created
	ErrUserNotCreated = errors.New("user not created")
)

// AccountService will be a placeholder
type AccountService interface {
	CreateAccount(context.Context, userEntities.User) error
}
