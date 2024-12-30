package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPaymentCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPaymentRepository(db)

	payment := &models.Payment{
		Reference:     "test-ref",
		Amount:        100,
		CampaignID:    "campaign-1",
		ContributorID: 1,
	}

	err := repo.Create(payment)
	assert.NoError(t, err)

	var saved models.Payment
	err = db.First(&saved, "reference = ?", payment.Reference).Error
	assert.NoError(t, err)
	assert.Equal(t, payment.Reference, saved.Reference)
}

func TestPayment_GetByReference(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPaymentRepository(db)

	payment := &models.Payment{
		Reference: "test-ref",
		Amount:    100,
	}
	db.Create(payment)

	found, err := repo.GetByReference(payment.Reference)
	assert.NoError(t, err)
	assert.Equal(t, payment.Reference, found.Reference)
}

func TestPayment_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPaymentRepository(db)

	payment := &models.Payment{
		Reference:     "test-ref",
		PaymentStatus: "pending",
	}
	db.Create(payment)

	payment.PaymentStatus = models.PaymentStatusSucceeded
	err := repo.Update(payment)
	assert.NoError(t, err)

	var updated models.Payment
	db.First(&updated, "reference = ?", payment.Reference)
	assert.Equal(t, models.PaymentStatusSucceeded, updated.PaymentStatus)
}

func TestPayment_GetByCampaign(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewPaymentRepository(db)

	// Create test payments
	payments := []models.Payment{
		{Reference: "ref1", CampaignID: "campaign-1"},
		{Reference: "ref2", CampaignID: "campaign-1"},
		{Reference: "ref3", CampaignID: "campaign-2"},
	}
	for _, p := range payments {
		db.Create(&p)
	}

	found, total, err := repo.GetByCampaign("campaign-1", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, found, 2)
}
