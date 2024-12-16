package contributor

import (
	"gorm.io/gorm"
)

type ContributorRepository interface {
	create(contribution *Contributor) error
	update(contribution *Contributor) error
	delete(contribution *Contributor) error
	processPayment(paymentID string) error
	refundPayment(paymentID string) error

	getContributors(campaignID string) ([]Contributor, error)
	getContributorById(contributorID uint) (Contributor, error)
	getContributorByUserHandle(userHandle uint) (Contributor, error)
	getEmailsOfExistingContributors(emails []string) ([]string, error)
}

type contributorRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) ContributorRepository {
	return &contributorRepository{db: db}
}

// ----------------------------------------------------------------------
func (r *contributorRepository) create(contribution *Contributor) error {
	return r.db.Create(contribution).Error
}
func (r *contributorRepository) update(contribution *Contributor) error {
	return r.db.Model(&Contributor{}).Where("id = ?", contribution.ID).Updates(map[string]interface{}{
		"amount":         contribution.Amount,
		"payment_status": contribution.PaymentStatus,
	}).Error
}

func (r *contributorRepository) delete(contribution *Contributor) error {
	return r.db.Delete(contribution).Error
}

func (r *contributorRepository) processPayment(paymentID string) error {
	// Implement payment processing logic here
	return nil
}

func (r *contributorRepository) refundPayment(paymentID string) error {
	// Implement payment refund logic here
	return nil
}

// ----------------------------------------------------------------------
func (r *contributorRepository) getContributors(campaignID string) ([]Contributor, error) {
	var contributors []Contributor
	err := r.db.Where("campaign_id = ?", campaignID).Find(&contributors).Error
	return contributors, err
}

func (r *contributorRepository) getContributorById(contributorID uint) (Contributor, error) {
	var contributor Contributor
	err := r.db.First(&contributor, contributorID).Error
	return contributor, err
}

func (r *contributorRepository) getContributorByUserHandle(userHandle uint) (Contributor, error) {
	var contributor Contributor
	err := r.db.Where("user_handle = ?", userHandle).First(&contributor).Error
	return contributor, err
}

func (r *contributorRepository) getEmailsOfExistingContributors(emails []string) ([]string, error) {
	var existingContributors []Contributor

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
