package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type PayoutRepository interface {
	Update(payout *models.Payout) error
	Create(payout *models.Payout) error
	GetByID(id string) (*models.Payout, error)
	GetByCampaignID(campaignID string, limit, offset int) ([]models.Payout, int64, error)
}
