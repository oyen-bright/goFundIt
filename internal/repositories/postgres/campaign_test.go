package postgress

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTestCampaign(db *gorm.DB, user models.User) models.Campaign {
	campaign := models.Campaign{
		ID:              "test-campaign-id-",
		Title:           "Test Campaign",
		CreatedBy:       user,
		CreatedByHandle: user.Handle,
		Description:     "Test Description of the campaign, Test Description of the campaign, Test Description of the campaign,Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign",
		TargetAmount:    1000,
		PaymentMethod:   models.PaymentMethodManual,
		Contributors: []models.Contributor{
			{

				CampaignID: "test-campaign-id-",
				Email:      "test@example.com",
				Amount:     1000,
			},
		},
		StartDate: time.Now(),

		EndDate: time.Now().Add(24 * time.Hour),
	}
	db.Create(&campaign)
	return campaign
}

func TestCampaignRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)

	campaign := &models.Campaign{
		ID:              "test-campaign-id",
		Title:           "Test Campaign",
		CreatedBy:       *user,
		CreatedByHandle: user.Handle,
		Description:     "Test Description of the campaign, Test Description of the campaign, Test Description of the campaign,Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign",
		TargetAmount:    1000,
		PaymentMethod:   models.PaymentMethodManual,
		Contributors: []models.Contributor{
			{

				CampaignID: "test-campaign-id",
				Email:      "test@example.com",
				Amount:     1000,
			},
		},
		StartDate: time.Now(),

		EndDate: time.Now().Add(24 * time.Hour),
	}

	result, err := repo.Create(campaign)
	assert.NoError(t, err)
	assert.Equal(t, campaign.ID, result.ID)
	assert.Equal(t, campaign.Title, result.Title)
}

func TestCampaignRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)
	campaign := createTestCampaign(db, *user)

	result, err := repo.GetByID(campaign.ID)
	assert.NoError(t, err)
	assert.Equal(t, campaign.ID, result.ID)
	assert.Equal(t, campaign.Title, result.Title)
}

func TestCampaignRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)
	campaign := createTestCampaign(db, *user)

	campaign.Title = "Updated Title"
	result, err := repo.Update(&campaign)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", result.Title)
}

func TestCampaignRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)
	campaign := createTestCampaign(db, *user)

	err = repo.Delete(campaign.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(campaign.ID)
	assert.Error(t, err) // Should return error as campaign is deleted
}

func TestCampaignRepository_GetByHandle(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)
	campaign := createTestCampaign(db, *user)

	result, err := repo.GetByHandle(user.Handle)
	assert.NoError(t, err)
	assert.Equal(t, campaign.ID, result.ID)
}

func TestCampaignRepository_GetExpiredCampaigns(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)

	// Create expired campaign
	expiredCampaign := models.Campaign{
		ID:              "test-campaign-id-",
		Title:           "Test Campaign",
		CreatedBy:       *user,
		CreatedByHandle: user.Handle,
		Description:     "Test Description of the campaign, Test Description of the campaign, Test Description of the campaign,Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign",
		TargetAmount:    1000,
		PaymentMethod:   models.PaymentMethodManual,
		Contributors: []models.Contributor{
			{

				CampaignID: "test-campaign-id",
				Email:      "test@example.com",
				Amount:     1000,
			},
		},
		StartDate: time.Now().Add(-48 * time.Hour),
		EndDate:   time.Now().Add(-24 * time.Hour), // Past date
	}
	err = db.Create(&expiredCampaign).Error
	assert.NoError(t, err)

	results, err := repo.GetExpiredCampaigns()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, expiredCampaign.ID, results[0].ID)
}

func TestCampaignRepository_GetActiveCampaigns(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)
	campaign := createTestCampaign(db, *user)

	results, err := repo.GetActiveCampaigns()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, campaign.ID, results[0].ID)
}

func TestCampaignRepository_GetNearEndCampaigns(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCampaignRepository(db)
	user, err := createTestUser(db)
	assert.NoError(t, err)

	// Create campaign ending soon
	nearEndCampaign := models.Campaign{
		ID:              "test-campaign-id-",
		Title:           "Test Campaign",
		CreatedBy:       *user,
		CreatedByHandle: user.Handle,
		Description:     "Test Description of the campaign, Test Description of the campaign, Test Description of the campaign,Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign Test Description of the campaign",
		TargetAmount:    1000,
		PaymentMethod:   models.PaymentMethodManual,
		Contributors: []models.Contributor{
			{

				CampaignID: "test-campaign-id",
				Email:      "test@example.com",
				Amount:     1000,
			},
		},
		StartDate: time.Now().Add(-78 * time.Hour),
		EndDate:   time.Now().Add(48 * time.Hour), // 2 days from now
	}
	db.Create(&nearEndCampaign)

	results, err := repo.GetNearEndCampaigns()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, nearEndCampaign.ID, results[0].ID)
}
