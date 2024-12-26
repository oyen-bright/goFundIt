package interfaces

import (
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/campaign"
	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignService interface {
	CreateCampaign(campaign *models.Campaign, userHandle string) (models.Campaign, error)
	UpdateCampaign(data dto.CampaignUpdateRequest, campaignID string, userHandle string) (*models.Campaign, error)
	DeleteCampaign(campaignID string) error

	GetCampaignByID(id string) (*models.Campaign, error)
	GetCampaignByIDWithContributors(id string) (*models.Campaign, error)
	GetCampaignByIDWithAllRelatedData(id string) (*models.Campaign, error)

	GetExpiredCampaigns() ([]models.Campaign, error)
	GetActiveCampaigns() ([]models.Campaign, error)
	GetNearEndCampaigns() ([]models.Campaign, error)

	RecalculateTargetAmount(campaignID string)
}
