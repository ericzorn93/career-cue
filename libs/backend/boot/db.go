package boot

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBOptions is a struct that defines the options for the database connection
type DBOptions struct {
	ConnectionDSN string
}

// IsZero checks if the DBOptions is empty
func (dbo DBOptions) IsZero() bool {
	return dbo.ConnectionDSN == ""
}

// BootServiceBuilder is a builder struct for the BootService Instance
func (bs *BootService) InitializeDB() (err error) {
	db, err := gorm.Open(postgres.Open(bs.dbOptions.ConnectionDSN), &gorm.Config{})
	if err != nil {
		bs.logger.Error("Failed to connect to database", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		bs.logger.Error("Failed to get generic database connection", err)
		return
	}

	bs.db = sqlDB
	bs.logger.Info("DB connected successfully")

	return err
}
