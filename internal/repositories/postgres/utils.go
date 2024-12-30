package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createTestUser(db *gorm.DB) (*models.User, error) {
	user := models.NewUser(
		"Test User",
		"test@example.com",
		true,
	)

	return user, db.Create(user).Error
}
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.Debug()

	err = db.AutoMigrate(
		&models.User{},
		&models.Contributor{},
		&models.Campaign{},
		&models.Otp{},
		&models.PlatformAnalytics{},
		&models.Comment{},
		&models.Activity{},
		&models.Payout{},
		&models.Activity{},
		&models.CampaignImage{},
		&models.Payment{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return db, cleanup
}
