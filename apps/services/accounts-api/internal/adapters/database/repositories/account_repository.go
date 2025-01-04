package repositories

import (
	"apps/services/accounts-api/internal/models"
	"context"
	"errors"
	"fmt"
	"libs/backend/boot"
	userEntities "libs/backend/domain/user/entities"
	userValueObjects "libs/backend/domain/user/valueobjects"
	"log/slog"
	"time"

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
		EmailAddress: user.EmailAddress.String(),
		UserName:     user.Username,
	}

	// Save the account in the database
	err := r.Database.Create(account).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAccountByCommonID gets an account from the database by commonID
func (r AccountRepository) GetAccountByCommonID(ctx context.Context, commonID userValueObjects.CommonID) (userEntities.User, error) {
	r.Logger.Info("Getting account", slog.String("commonID", commonID.String()))

	account := &models.Account{}
	r.Database.First(account, "common_id = ?", commonID.Value())

	if account.ID == uuid.Nil {
		return userEntities.User{}, errors.New("account not found")
	}

	return r.convertAccountToUser(account), nil
}

// GetAccountByEmailAddress gets an account from the database by email address
func (r AccountRepository) GetAccountByEmailAddress(ctx context.Context, emailAdress userValueObjects.EmailAddress) (userEntities.User, error) {
	r.Logger.Info("Getting account", slog.String("emailAdress", emailAdress.String()))

	account := &models.Account{}
	r.Database.First(account, "email_address = ?", emailAdress)

	if account.ID == uuid.Nil {
		return userEntities.User{}, errors.New("account not found")
	}

	return r.convertAccountToUser(account), nil
}

// SoftDeleteAccountByCommonID will mark the user as deleted in the database with a timestamp
func (r AccountRepository) SoftDeleteAccountByCommonID(ctx context.Context, commonID userValueObjects.CommonID) (time.Time, error) {
	r.Logger.Info("Handling soft deletion of account by commonID", slog.String("commonID", commonID.String()))

	// Handle soft deletion in the database
	deletedAccount := &models.Account{}
	if err := r.Database.Delete(deletedAccount, "common_id = ?", commonID).Error; err != nil {
		r.Logger.Error("Cannot soft delete the user by commonID", slog.String("commonID", commonID.String()))
		return time.Time{}, fmt.Errorf("cannot soft delete account: %w", err)
	}

	return deletedAccount.DeletedAt.Time, nil
}

// HardDeleteAccountByCommonID will remove the user from the dataase
func (r AccountRepository) HardDeleteAccountByCommonID(ctx context.Context, commonID userValueObjects.CommonID) (time.Time, error) {
	r.Logger.Info("Handling hard deletion of account by commonID", slog.Any("commonID", commonID.String()))

	// Handle soft deletion in the database
	deletedAccount := &models.Account{}
	if err := r.Database.Unscoped().Delete(deletedAccount, "common_id = ?", commonID).Error; err != nil {
		r.Logger.Error("Cannot hard delete the user by commonID", slog.String("commonID", commonID.String()))
		return time.Time{}, fmt.Errorf("cannot hard delete accoutn: %w", err)
	}

	return deletedAccount.DeletedAt.Time, nil
}

// convertAccountToUser converts an account to a user
func (r AccountRepository) convertAccountToUser(account *models.Account) userEntities.User {
	parsedCommonID := userValueObjects.NewCommonIDFromUUID(account.CommonID)
	parsedEmailAddress := userValueObjects.NewEmailAddress(account.EmailAddress)

	return userEntities.NewUser(
		userEntities.WithCommonID(parsedCommonID),
		userEntities.WithEmailAddress(parsedEmailAddress),
		userEntities.WithUserUsername(account.UserName),
		userEntities.WithCreatedAt(account.CreatedAt),
		userEntities.WithUpdatedAt(account.UpdatedAt),
	)
}
