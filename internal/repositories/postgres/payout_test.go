package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPayoutCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPayoutRepository(db)

	payout := models.NewPayout("campaign-1", 1000, models.PaymentMethodManual)

	err := repo.Create(payout)
	assert.NoError(t, err)

	var saved models.Payout
	err = db.First(&saved, "id = ?", payout.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, payout.ID, saved.ID)
}

func TestPayoutUpdate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPayoutRepository(db)

	payout := models.NewPayout("campaign-1", 1000, models.PaymentMethodManual)

	db.Create(payout)

	payout.Status = models.PayoutStatusCompleted
	err := repo.Update(payout)
	assert.NoError(t, err)

	var updated models.Payout
	db.First(&updated, "id = ?", payout.ID)
	assert.Equal(t, models.PayoutStatusCompleted, updated.Status)
}

func TestPayoutGetByCampaignID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPayoutRepository(db)

	payouts := []models.Payout{
		*models.NewPayout("1", 1000, models.PaymentMethodManual),
		*models.NewPayout("2", 1000, models.PaymentMethodManual),
		*models.NewPayout("1", 1000, models.PaymentMethodManual),
	}

	for _, p := range payouts {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	var allPayouts []models.Payout
	err := db.Find(&allPayouts).Error
	assert.NoError(t, err)

	result, _, err := repo.GetByCampaignID("1", 3, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].CampaignID)
	assert.Equal(t, "1", result[1].CampaignID)
}

func TestPayoutGetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPayoutRepository(db)

	payout := models.NewPayout("campaign-1", 1000, models.PaymentMethodManual)
	err := db.Create(payout).Error
	assert.NoError(t, err)

	found, err := repo.GetByID(payout.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, payout.ID, found.ID)
	assert.Equal(t, payout.CampaignID, found.CampaignID)
	assert.Equal(t, payout.Amount, found.Amount)

	notFound, err := repo.GetByID("non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, notFound)
}
