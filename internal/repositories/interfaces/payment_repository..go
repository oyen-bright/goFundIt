package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type PaymentRepository interface {
	CreatePayment(payment *models.Payment) error
	GetPayment(paymentID string) (*models.Payment, error)
	GetPaymentByContributorID(contributorID string) ([]*models.Payment, error)
	UpdatePayment(payment *models.Payment) error
	DeletePayment(paymentID string) error
	ListPayments(campaignID string) ([]*models.Payment, error)
}
