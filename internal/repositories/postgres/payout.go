package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type payoutRepository struct {
	db *gorm.DB
}

// Create implements interfaces.PayoutRepository.
func (p *payoutRepository) Create(payout *models.Payout) error {
	return p.db.Create(payout).Error
}

// GetByCampaignID implements interfaces.PayoutRepository.
func (p *payoutRepository) GetByCampaignID(campaignID string, limit int, offset int) ([]models.Payout, int64, error) {
	var payouts []models.Payout
	var total int64

	query := p.db.Model(&models.Payout{}).Where("campaign_id = ?", campaignID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&payouts).Error; err != nil {
		return nil, 0, err
	}

	return payouts, total, nil
}

// GetByID implements interfaces.PayoutRepository.
func (p *payoutRepository) GetByID(id string) (*models.Payout, error) {
	var payout models.Payout
	if err := p.db.First(&payout, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payout, nil
}

// Update implements interfaces.PayoutRepository.
func (p *payoutRepository) Update(payout *models.Payout) error {
	return p.db.Save(payout).Error
}

// NewPayoutRepository creates a new instance of the payout repository
func NewPayoutRepository(db *gorm.DB) interfaces.PayoutRepository {
	return &payoutRepository{db: db}
}
