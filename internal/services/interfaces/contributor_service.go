package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ContributorService interface {
	AddContributorToCampaign(contribution *models.Contributor, campaignId, userHandle string) error
	UpdateContributor(contributor *models.Contributor) error
	UpdateContributorByID(contributor *models.Contributor, contributorID uint, userEmail string) (retrievedContributor models.Contributor, err error)
	GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error)
	RemoveContributorFromCampaign(contributorId uint, campaignId, userHandle string) error
	GetContributorByID(contributorID uint) (models.Contributor, error)
	GetContributorByIDWithActivities(contributorID uint) (models.Contributor, error)
	// GetEmailsOfActiveContributors(emails []string) ([]string, error)
	// GetContributorByUserHandle(userHandle uint) (models.Contributor, error)
	// ContributeToCampaign(userID uint, campaignID string, amount float64) error
	// ProcessPayment(paymentID string) error
	// RefundPayment(paymentID string) error
}
