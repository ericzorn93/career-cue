package services

import (
	"apps/services/accounts-api/internal/app/ports"
	"context"
	"errors"
	"libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
	"log/slog"
)

// RegistrationService is the registration service
type RegistrationService struct {
	// Logger is the logger from the boot framework
	Logger boot.Logger

	// AccountRepository is the account repository
	AccountRepository ports.AccountRepository
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(logger boot.Logger, accountRepository ports.AccountRepository) RegistrationService {
	return RegistrationService{
		Logger:            logger,
		AccountRepository: accountRepository,
	}
}

// RegisterUser registers a user in the system and the database
func (s RegistrationService) RegisterUser(ctx context.Context, user userEntities.User) error {
	s.Logger.Info("Registering user", slog.String("commonID", user.CommonID.String()))

	// Create account
	err := s.AccountRepository.CreateAccount(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the system
func (s RegistrationService) GetUser(ctx context.Context, commonID userValueObjects.CommonID, emailAddress userValueObjects.EmailAddress) (userEntities.User, error) {
	var err error

	s.Logger.Info("Getting user", slog.String("commonID", commonID.String()))

	// Get account from the repository by commonID or email address
	var user userEntities.User

	switch {
	case !commonID.IsEmpty():
		user, err = s.AccountRepository.GetAccountByCommonID(ctx, commonID)
		if err != nil {
			return userEntities.NewEmptyUser(), err
		}
	case !emailAddress.IsEmpty():
		user, err = s.AccountRepository.GetAccountByEmailAddress(ctx, emailAddress)
		if err != nil {
			return userEntities.NewEmptyUser(), err
		}
	default:
		return userEntities.NewEmptyUser(), errors.New("commonID or emailAddress must be provided")
	}

	return user, nil
}
