package repositories

import (
	"apps/services/accounts-api/internal/models"
	"context"
	"errors"
	"libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
	"log/slog"

	"github.com/google/uuid"
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

// GetAccount gets an account from the database by commonID
func (r AccountRepository) GetAccount(ctx context.Context, commonID userValueObjects.CommonID) (userEntities.User, error) {
	r.Logger.Info("Getting account", slog.String("commonID", commonID.String()))

	account := &models.Account{}
	r.Database.First(account, "common_id = ?", commonID.Value())

	if account.ID == uuid.Nil {
		return userEntities.User{}, errors.New("account not found")
	}

	parsedCommonID := userValueObjects.NewCommonIDFromUUID(account.CommonID)

	return userEntities.NewUser(
		userEntities.WithCommonID(parsedCommonID),
		userEntities.WithEmailAddress(account.EmailAddress),
		userEntities.WithUserUsername(account.UserName),
		userEntities.WithCreatedAt(account.CreatedAt),
		userEntities.WithUpdatedAt(account.UpdatedAt),
	), nil
}
