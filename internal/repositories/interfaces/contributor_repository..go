package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ContributorRepository interface {
	Create(contribution *models.Contributor) error
	Update(contribution *models.Contributor) error
	UpdateName(contributorID uint, name string) error
	Delete(contribution *models.Contributor) error
	ProcessPayment(paymentID string) error
	RefundPayment(paymentID string) error

	GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error)
	GetContributorById(contributorID uint) (models.Contributor, error)
	GetContributorByUserHandle(userHandle uint) (models.Contributor, error)
	GetEmailsOfActiveContributors(emails []string) ([]string, error)
}
