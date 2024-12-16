package migrations

import (
	"github.com/oyen-bright/goFundIt/internal/activity"
	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/campaign"
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/payment"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Applying migration
	// DropOtpTable(db)
	err := db.AutoMigrate(
		&otp.Otp{},
		&campaign.Campaign{},
		&campaign.CampaignImage{},
		&activity.Activity{},
		&contributor.Contributor{},
		&auth.User{},
		&payment.Payment{},
	)
	if err != nil {
		return err
	}

	return nil
}

func DropOtpTable(db *gorm.DB) error {
	// Dropping the otp table
	err := db.Migrator().DropTable(&campaign.Campaign{}, &campaign.CampaignImage{}, auth.User{}, contributor.Contributor{}, otp.Otp{}, activity.Activity{})
	if err != nil {
		return err
	}
	return nil
}
