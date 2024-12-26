package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
)

type PaymentService interface {
	InitializePayment(contributorID uint) (*models.Payment, error)
	InitializeManualPayment(contributorID uint, reference, userEmail string) (*models.Payment, error)

	VerifyPayment(reference string) error
	VerifyManualPayment(reference, userHandle string) error

	DeletePayment(payment models.Payment) error

	GetPaymentByReference(reference string) (*models.Payment, error)
	GetPaymentsByCampaign(campaignID string, limit, offset int) ([]*models.Payment, int64, error)
	GetPaymentsByContributor(contributorID uint, limit, offset int) ([]models.Payment, int64, error)

	ProcessPaystackWebhook(event paystack.PaystackWebhookEvent)
}
