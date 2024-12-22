package migrations

import (
	"github.com/oyen-bright/goFundIt/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Applying migration
	// DropOtpTable(db)
	err := db.AutoMigrate(
		&models.Otp{},
		&models.Campaign{},
		&models.CampaignImage{},
		&models.Activity{},
		&models.Contributor{},
		&models.User{},
		&models.Payment{},
		&models.Comment{},
		&models.Payout{},
	)
	if err != nil {
		return err
	}

	return nil
}

func DropOtpTable(db *gorm.DB) error {
	// Dropping the otp table
	err := db.Migrator().DropTable(&models.Campaign{}, &models.CampaignImage{}, models.User{}, models.Contributor{}, models.Otp{}, models.Activity{})
	if err != nil {
		return err
	}
	return nil
}
