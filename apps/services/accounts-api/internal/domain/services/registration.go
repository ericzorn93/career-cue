package services

import (
	"apps/services/accounts-api/internal/app/ports"
	"context"
	"libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
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
