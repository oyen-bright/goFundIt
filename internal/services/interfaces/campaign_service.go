package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignService interface {
	CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error)
	GetCampaignByID(id string) (*models.Campaign, error)
	GetCampaignByIDWithContributors(id string) (*models.Campaign, error)

	CheckExistingCampaign(userHandle string) error
	// UpdateCampaign(campaign Campaign) (Campaign, error)
	// DeleteCampaign(campaignID string) error
	// JoinCampaign(userID uint, campaignID string) error
	// LeaveCampaign(userID uint, campaignID string) error
	// SetCampaignTheme(campaignID string, themeID uint) error
	// GetCampaignsByUser(handle string) ([]Campaign, error)
	// UserCanContribute(handle string, campaignID string) error
}
