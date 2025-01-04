package database

import (
	"fmt"

	"github.com/oyen-bright/goFundIt/pkg/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

func Init(cfg Config, isDevelopment bool) (*gorm.DB, error) {
	// Setup PostgreSQL connection using the provided configuration
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)
	logMode := logger.Warn
	if isDevelopment {
		logMode = logger.Info
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		return nil, err
	}

	err = migrations.Migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil

}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
