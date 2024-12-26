package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
	Delete(reference string) error

	GetByReference(reference string) (*models.Payment, error)
	GetByContributor(contributorID uint, limit, offset int) ([]models.Payment, int64, error)
	GetByCampaign(campaignID string, limit, offset int) ([]*models.Payment, int64, error)
}
