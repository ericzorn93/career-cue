package boot

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBOptions is a struct that defines the options for the database connection
type DBOptions struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     string
	SSLMode  string
	TimeZone string
}

// IsZero checks if the DBOptions is empty
func (dbo DBOptions) IsZero() bool {
	return dbo.Host == "" &&
		dbo.Name == "" &&
		dbo.User == "" &&
		dbo.Password == "" &&
		dbo.Port == "" &&
		dbo.SSLMode == "" &&
		dbo.TimeZone == ""

}

// BootServiceBuilder is a builder struct for the BootService Instance
func (bs *BootService) InitializeDB() (err error) {
	// Construct the DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", bs.dbOptions.Host, bs.dbOptions.User, bs.dbOptions.Password, bs.dbOptions.Name, bs.dbOptions.Port, bs.dbOptions.SSLMode, bs.dbOptions.TimeZone)
	bs.logger.Info("Connecting to database", slog.String("dsn", dsn))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		bs.logger.Error("Failed to connect to database", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		bs.logger.Error("Failed to get generic database connection", err)
		return
	}

	bs.localDB = sqlDB
	bs.db = db
	bs.logger.Info("DB connected successfully")

	return err
}
