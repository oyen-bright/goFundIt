package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignRepository interface {
	Create(campaign *models.Campaign) (models.Campaign, error)
	Update(campaign *models.Campaign) (models.Campaign, error)
	Delete(campaignID string) error

	GetByID(id string, preload bool) (models.Campaign, error)
	GetByIDWithContributors(id string) (models.Campaign, error)
	GetByCreatorHandle(handle string, preload bool) (models.Campaign, error)
}
