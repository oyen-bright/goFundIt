package database

import (
	"fmt"

	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/internal/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg config.AppConfig) (*gorm.DB, error) {
	// Setup PostgreSQL connection using the provided configuration
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	DB = db

	err = migrations.Migrate(DB)
	if err != nil {
		return nil, err
	}

	return DB, nil

}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
