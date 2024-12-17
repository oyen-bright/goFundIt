package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type CampaignService interface {
	CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error)
	GetCampaignByID(id string) (*models.Campaign, error)

	UserCanCreateCampaign(userHandle string) error
	EmailsCanContribute(contributorsEmail []string) ([]string, error)
	// UpdateCampaign(campaign Campaign) (Campaign, error)
	// DeleteCampaign(campaignID string) error
	// JoinCampaign(userID uint, campaignID string) error
	// LeaveCampaign(userID uint, campaignID string) error
	// SetCampaignTheme(campaignID string, themeID uint) error
	// GetCampaignsByUser(handle string) ([]Campaign, error)
	// UserCanContribute(handle string, campaignID string) error
}
