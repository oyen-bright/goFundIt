package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ContributorService interface {
	AddContributorToCampaign(contribution *models.Contributor) error
	RemoveContributorFromCampaign(contribution *models.Contributor) error
	CanContributeToCampaign(userID uint, campaignID string) (bool, error)
	GetContributors(campaignID string) ([]models.Contributor, error)
	GetContributorByID(contributorID uint) (models.Contributor, error)
	GetContributorByUserHandle(userHandle uint) (models.Contributor, error)
	GetEmailsOfExistingContributors(emails []string) ([]string, error)
	ContributeToCampaign(userID uint, campaignID string, amount float64) error
	ProcessPayment(paymentID string) error
	RefundPayment(paymentID string) error
}
