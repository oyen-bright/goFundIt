package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignRepository interface {
	Create(campaign *models.Campaign) (models.Campaign, error)
	Update(campaign *models.Campaign) (models.Campaign, error)
	Delete(campaignID string) error

	GetByID(id string) (models.Campaign, error)
	GetByIDWithSelectedData(id string, options models.PreloadOption) (models.Campaign, error)
	GetByHandle(handle string) (models.Campaign, error)

	GetExpiredCampaigns() ([]models.Campaign, error)
	GetActiveCampaigns() ([]models.Campaign, error)
	GetNearEndCampaigns() ([]models.Campaign, error)
}
