package interfaces

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignService interface {
	CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error)
	GetCampaignByID(id string) (*models.Campaign, error)
	GetCampaignByIDWithContributors(id string) (*models.Campaign, error)
	GetCampaignByIDWithAllRelatedData(id string) (*models.Campaign, error)
	DeleteCampaign(campaignID string) error

	GetExpiredCampaigns() ([]models.Campaign, error)
	GetActiveCampaigns() ([]models.Campaign, error)
	GetNearEndCampaigns() ([]models.Campaign, error)

	GetCampaignsForAnalytics(yesterday time.Time, today time.Time) ([]models.Campaign, error)
}
