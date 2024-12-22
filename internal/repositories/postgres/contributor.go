package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type contributorRepository struct {
	db *gorm.DB
}

func NewContributorRepository(db *gorm.DB) interfaces.ContributorRepository {
	return &contributorRepository{db: db}
}

// ----------------------------------------------------------------------
func (r *contributorRepository) Create(contribution *models.Contributor) error {
	return r.db.Create(contribution).Error
}

// UpdateName updates only the name field for a contributor
func (r *contributorRepository) UpdateName(contributorID uint, name string) error {
	return r.db.Model(&models.Contributor{}).
		Where("id = ?", contributorID).
		Update("name", name).Error
}

func (r *contributorRepository) Update(contribution *models.Contributor) error {
	//TODO:use gorm struct tag to annotate the fields to be updated
	return r.db.Model(&models.Contributor{}).Where("id = ?", contribution.ID).Updates(map[string]interface{}{
		"amount":         contribution.Amount,
		"name":           contribution.Name,
		"amount_paid":    contribution.AmountPaid,
		"payment_status": contribution.PaymentStatus,
	}).Error
}

func (r *contributorRepository) Delete(contribution *models.Contributor) error {
	return r.db.Delete(contribution).Error
}

func (r *contributorRepository) ProcessPayment(paymentID string) error {
	// Implement payment processing logic here
	return nil
}

func (r *contributorRepository) RefundPayment(paymentID string) error {
	// Implement payment refund logic here
	return nil
}

// ----------------------------------------------------------------------
func (r *contributorRepository) GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error) {
	var contributors []models.Contributor
	err := r.db.Where("campaign_id = ?", campaignID).Find(&contributors).Error
	return contributors, err
}

func (r *contributorRepository) GetContributorById(contributorID uint, preload bool) (models.Contributor, error) {
	var contributor models.Contributor

	if preload {
		var contributor models.Contributor
		err := r.db.Preload("Activities").First(&contributor, contributorID).Error
		return contributor, err
	}
	err := r.db.First(&contributor, contributorID).Error
	return contributor, err
}

func (r *contributorRepository) GetContributorByUserHandle(userHandle uint) (models.Contributor, error) {
	var contributor models.Contributor
	err := r.db.Where("user_handle = ?", userHandle).First(&contributor).Error
	return contributor, err
}

func (r *contributorRepository) CanContributeToCampaign(userID uint, campaignID string) (bool, error) {
	var contributor models.Contributor
	err := r.db.Where("email = ? AND campaign_id = ?", userID, campaignID).First(&contributor).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *contributorRepository) GetEmailsOfActiveContributors(emails []string) ([]string, error) {
	var existingContributors []models.Contributor

	err := r.db.Where("email IN (?)", emails).Find(&existingContributors).Error
	if err != nil {
		return nil, err
	}

	existingEmails := make([]string, 0, len(existingContributors))
	for _, contributor := range existingContributors {
		existingEmails = append(existingEmails, contributor.Email)
	}

	return existingEmails, nil
}
