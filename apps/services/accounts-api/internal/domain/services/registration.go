package services

import (
	"apps/services/accounts-api/internal/app/ports"
	"context"
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
func (s RegistrationService) GetUser(ctx context.Context, commonID userValueObjects.CommonID) (userEntities.User, error) {
	s.Logger.Info("Getting user", slog.String("commonID", commonID.String()))

	// Get account from the repository
	user, err := s.AccountRepository.GetAccount(ctx, commonID)
	if err != nil {
		return userEntities.NewEmptyUser(), err
	}

	return user, nil
}
