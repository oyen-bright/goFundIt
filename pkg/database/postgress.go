package database

import (
	"fmt"

	"github.com/oyen-bright/goFundIt/pkg/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

func Init(cfg Config) (*gorm.DB, error) {
	// Setup PostgreSQL connection using the provided configuration
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.Logger = db.Logger.LogMode(3)

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
