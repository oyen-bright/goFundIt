package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepository{db: db}
}
func (r *paymentRepository) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetPayment(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Preload("Contributor").First(&payment, "id = ?", paymentID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) DeletePayment(paymentID string) error {
	return r.db.Delete(&models.Payment{}, "id = ?", paymentID).Error
}

func (r *paymentRepository) ListPayments(campaignID string) ([]*models.Payment, error) {
	var payments []*models.Payment
	if err := r.db.Where("campaign_id = ?", campaignID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) GetPaymentByContributorID(contributorID string) ([]*models.Payment, error) {
	var payments []*models.Payment
	if err := r.db.Where("contributor_id = ?", contributorID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
