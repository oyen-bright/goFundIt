package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type ContributorRepository struct {
	db *gorm.DB
}

func NewContributorRepository(db *gorm.DB) interfaces.ContributorRepository {
	return &ContributorRepository{db: db}
}

// ----------------------------------------------------------------------
func (r *ContributorRepository) Create(contribution *models.Contributor) error {
	return r.db.Create(contribution).Error
}

func (r *ContributorRepository) Update(contribution *models.Contributor) error {
	return r.db.Model(&models.Contributor{}).Where("id = ?", contribution.ID).Updates(map[string]interface{}{
		"amount":         contribution.Amount,
		"payment_status": contribution.PaymentStatus,
	}).Error
}

func (r *ContributorRepository) Delete(contribution *models.Contributor) error {
	return r.db.Delete(contribution).Error
}

func (r *ContributorRepository) ProcessPayment(paymentID string) error {
	// Implement payment processing logic here
	return nil
}

func (r *ContributorRepository) RefundPayment(paymentID string) error {
	// Implement payment refund logic here
	return nil
}

// ----------------------------------------------------------------------
func (r *ContributorRepository) GetContributors(campaignID string) ([]models.Contributor, error) {
	var contributors []models.Contributor
	err := r.db.Where("campaign_id = ?", campaignID).Find(&contributors).Error
	return contributors, err
}

func (r *ContributorRepository) GetContributorById(contributorID uint) (models.Contributor, error) {
	var contributor models.Contributor
	err := r.db.First(&contributor, contributorID).Error
	return contributor, err
}

func (r *ContributorRepository) GetContributorByUserHandle(userHandle uint) (models.Contributor, error) {
	var contributor models.Contributor
	err := r.db.Where("user_handle = ?", userHandle).First(&contributor).Error
	return contributor, err
}

func (r *ContributorRepository) GetEmailsOfExistingContributors(emails []string) ([]string, error) {
	var existingContributors []models.Contributor

	err := r.db.Where("user_email IN (?)", emails).Find(&existingContributors).Error
	if err != nil {
		return nil, err
	}

	existingEmails := make([]string, 0, len(existingContributors))
	for _, contributor := range existingContributors {
		existingEmails = append(existingEmails, contributor.UserEmail)
	}

	return existingEmails, nil
}
