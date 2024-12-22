package services

import (
	"fmt"
	"net/http"

	"github.com/oyen-bright/goFundIt/internal/models"
	repos "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
)

type paymentService struct {
	repo               repos.PaymentRepository
	paystack           *paystack.Client
	campaignService    services.CampaignService
	contributorService services.ContributorService
	logger             logger.Logger
}

// VerifyPayment implements interfaces.PaymentService.
func (p *paymentService) VerifyPayment(reference string) error {
	// Get the payment
	payment, err := p.repo.GetByReference(reference)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.NotFound("Payment not found")
		}
		return errs.InternalServerError(err).Log(p.logger)
	}
	// Verify the payment
	res, err := p.paystack.VerifyTransaction(reference)
	if err != nil {
		return errs.InternalServerError(err).Log(p.logger)
	}

	// Check if the payment is successful
	if res.IsPaymentSuccessful() {
		// Update the payment status

		payment.SetPaymentStatusToSuccess()
		gatewayResponse := res.ToString()
		payment.GatewayResponse = &gatewayResponse
		err = p.repo.Update(payment)
		if err != nil {
			return errs.InternalServerError(err).Log(p.logger)
		}

		// Update the contributor status
		contributor := payment.Contributor
		//TODO: amount paid "amount": 2000,   "amountPaid": 210000, should be 2100
		contributor.SetPaymentSucceeded(float64(res.Data.Amount))
		err = p.contributorService.UpdateContributor(&contributor)
		if err != nil {
			return errs.InternalServerError(err).Log(p.logger)
		}
		return nil

	}

	//TODO: Send email to contributor
	//TODO: set payment status to failed

	return errs.New(fmt.Sprintf("Payment Verification failed :%v ", res.Data.GatewayResponse), http.StatusUnprocessableEntity)

}

// HandlePaystackWebhook implements interfaces.PaymentService.
func (p *paymentService) HandlePaystackWebhook(event paystack.PaystackWebhookEvent) {
	// validate payment
	payment, err := p.repo.GetByReference(event.Data.Reference)
	if err != nil {
		errs.InternalServerError(err).Log(p.logger)
		return
	}
	contributor := payment.Contributor
	// validate event type
	switch event.Event {
	// Handle the charge success event
	case paystack.EventChargeSuccess:
		contributor.SetPaymentSucceeded(event.Data.Amount)
		payment.SetPaymentStatusToSuccess()
		err = p.contributorService.UpdateContributor(&contributor)
		if err != nil {
			errs.InternalServerError(err).Log(p.logger)
		}
		return
		// Handle the charge failed event
	case paystack.EventChargeFailed:
		contributor.SetPaymentFailed()
		payment.SetPaymentStatusToFailed()
		err = p.contributorService.UpdateContributor(&contributor)
		if err != nil {
			errs.InternalServerError(err).Log(p.logger)
			return
		}
		return

	}
	// update the payment status
	err = p.repo.Update(payment)
	if err != nil {
		errs.InternalServerError(err).Log(p.logger)
	}
}

// InitializePayment implements interfaces.PaymentService.
func (p *paymentService) InitializePayment(contributorID uint) (*models.Payment, error) {

	// Validate the contributor
	contributor, err := p.contributorService.GetContributorByIDWithActivities(contributorID)
	if err != nil {
		return nil, err
	}
	if contributor.HasPaid() {
		return nil, errs.BadRequest("Contributor has already paid", nil)
	}

	// validate campaign
	campaign, err := p.campaignService.GetCampaignByID(contributor.CampaignID)

	if err != nil {
		return nil, err
	}
	if campaign.HasEnded() {
		return nil, errs.BadRequest("Campaign has ended", nil)
	}

	if campaign.PaymentMethod == models.PaymentMethodManual {
		return nil, errs.BadRequest("Campaign payment method is manual", nil)
	}

	// Create a payment
	if campaign.PaymentMethod == models.PaymentMethodFiat {
		response, err := p.paystack.InitiateTransaction(contributor.Email, string(*campaign.FiatCurrency), contributor.GetAmountTotal())
		if err != nil {
			return nil, errs.InternalServerError(err).Log(p.logger)
		}
		payment := models.NewFiatPayment(contributor.ID, campaign.ID, response.Data.Reference, contributor.GetAmountTotal(), response.Data.AuthorizationURL)
		// Save the payment

		err = p.repo.Create(payment)

		if err != nil {
			return nil, errs.InternalServerError(err).Log(p.logger)
		}
		return payment, nil
	}

	return nil, errs.BadRequest("Invalid payment method", nil)

}

// DeletePayment implements interfaces.PaymentService.
func (p *paymentService) DeletePayment(payment models.Payment) error {
	panic("unimplemented")
}

// GetPaymentByReference implements interfaces.PaymentService.
func (p *paymentService) GetPaymentByReference(reference string) (*models.Payment, error) {
	panic("unimplemented")
}

// GetPaymentsByCampaign implements interfaces.PaymentService.
func (p *paymentService) GetPaymentsByCampaign(campaignID string, limit int, offset int) ([]*models.Payment, int64, error) {
	panic("unimplemented")
}

// GetPaymentsByContributor implements interfaces.PaymentService.
func (p *paymentService) GetPaymentsByContributor(contributorID uint, limit int, offset int) ([]models.Payment, int64, error) {
	panic("unimplemented")
}
func NewPaymentService(repo repos.PaymentRepository, contributorService services.ContributorService, campaignService services.CampaignService, paystack *paystack.Client, logger logger.Logger) services.PaymentService {
	return &paymentService{
		repo:               repo,
		campaignService:    campaignService,
		logger:             logger,
		paystack:           paystack,
		contributorService: contributorService,
	}
}
