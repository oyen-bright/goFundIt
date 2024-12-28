package migrations

import (
	"github.com/oyen-bright/goFundIt/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Applying migration
	// DropOtpTable(db)
	err := db.AutoMigrate(
		// Base tables (no foreign key dependencies)
		&models.User{},
		&models.Otp{},

		&models.PlatformAnalytics{},

		&models.Campaign{},
		&models.CampaignImage{},

		&models.Payout{},
		&models.Contributor{},
		&models.Comment{},
		&models.Activity{},
		&models.Payment{},
	)
	if err != nil {
		return err
	}

	return nil
}

func DropOtpTable(db *gorm.DB) error {
	// Dropping the otp table
	err := db.Migrator().DropTable(
		// &models.User{},
		// &models.Otp{},
		// &models.PlatformAnalytics{},

		// &models.Campaign{},
		// &models.CampaignImage{},

		// &models.Payout{},
		// &models.Contributor{},
		// &models.Comment{},
		// &models.Activity{},
		&models.Payment{},
	)
	if err != nil {
		return err
	}
	return nil
}
