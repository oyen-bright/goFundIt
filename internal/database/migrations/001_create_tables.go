package migrations

import (
	"github.com/oyen-bright/goFundIt/internal/otp/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Applying migration
	err := db.AutoMigrate(
		&model.Otp{},
	)
	if err != nil {
		return err
	}
	return nil
}
func ClearOtpTable(db *gorm.DB) error {
	// Deleting all data in the otp table
	err := db.Exec("DELETE FROM otps").Error
	if err != nil {
		return err
	}
	return nil
}
func DropOtpTable(db *gorm.DB) error {
	// Dropping the otp table
	err := db.Migrator().DropTable(&model.Otp{})
	if err != nil {
		return err
	}
	return nil
}
