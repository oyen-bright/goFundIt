package payment

import "gorm.io/gorm"

type PaymentRepository interface {
	CreatePayment(payment *Payment) error
	GetPayment(paymentID string) (*Payment, error)
	GetPaymentByContributorID(contributorID string) ([]*Payment, error)
	UpdatePayment(payment *Payment) error
	DeletePayment(paymentID string) error
	ListPayments(campaignID string) ([]*Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}
func (r *paymentRepository) CreatePayment(payment *Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetPayment(paymentID string) (*Payment, error) {
	var payment Payment
	if err := r.db.Preload("Contributor").First(&payment, "id = ?", paymentID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePayment(payment *Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) DeletePayment(paymentID string) error {
	return r.db.Delete(&Payment{}, "id = ?", paymentID).Error
}

func (r *paymentRepository) ListPayments(campaignID string) ([]*Payment, error) {
	var payments []*Payment
	if err := r.db.Where("campaign_id = ?", campaignID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) GetPaymentByContributorID(contributorID string) ([]*Payment, error) {
	var payments []*Payment
	if err := r.db.Where("contributor_id = ?", contributorID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
