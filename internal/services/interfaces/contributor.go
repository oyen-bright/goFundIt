package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ContributorService interface {
	GetContributorByID(contributorID uint) (models.Contributor, error)

	UpdateContributor(contributor *models.Contributor) error
	UpdateContributorByID(contributor *models.Contributor, contributorID uint, userEmail string) (retrievedContributor models.Contributor, err error)

	GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error)

	AddContributorToCampaign(contribution *models.Contributor, campaignId, campaignKey, userHandle string) error
	RemoveContributorFromCampaign(contributorId uint, campaignId, userHandle, key string) error
}
