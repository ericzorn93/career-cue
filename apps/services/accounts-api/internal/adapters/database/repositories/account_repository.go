package repositories

import (
	"apps/services/accounts-api/internal/models"
	"context"
	"libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	"log/slog"

	"gorm.io/gorm"
)

// AccountRepository is the account repository that interacts with the database
type AccountRepository struct {
	// Logger is the logger from the boot framework
	Logger boot.Logger

	// Database is the database connection
	Database *gorm.DB
}

// NewAccountRespository creates a new account repository
func NewAccountRespository(logger boot.Logger, db *gorm.DB) AccountRepository {
	return AccountRepository{
		Logger:   logger,
		Database: db,
	}
}

// CreateAccount creates an account in the database
func (r AccountRepository) CreateAccount(ctx context.Context, user userEntities.User) error {
	r.Logger.Info("Creating account", slog.String("commonID", user.CommonID.String()))

	// Create account
	account := &models.Account{
		CommonID:     user.CommonID.Value(),
		EmailAddress: user.EmailAddress,
		UserName:     user.Username,
	}

	// Save the account in the database
	err := r.Database.Create(account).Error
	if err != nil {
		return err
	}

	return nil
}
