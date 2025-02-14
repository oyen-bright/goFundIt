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
	"github.com/oyen-bright/goFundIt/pkg/storage"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

type paymentService struct {
	repo                repos.PaymentRepository
	paystack            paystack.PaystackClient
	campaignService     services.CampaignService
	analyticsService    services.AnalyticsService
	contributorService  services.ContributorService
	notificationService services.NotificationService
	broadcaster         services.EventBroadcaster
	storage             storage.Storage
	logger              logger.Logger
	runAsync            func(func())
}

func NewPaymentService(
	repo repos.PaymentRepository,
	contributorService services.ContributorService,
	analyticsService services.AnalyticsService,
	campaignService services.CampaignService,
	notificationService services.NotificationService,
	paystack paystack.PaystackClient,
	storage storage.Storage,
	broadcaster services.EventBroadcaster,
	logger logger.Logger,
) services.PaymentService {
	return &paymentService{
		// Repository
		repo: repo,

		// Services
		campaignService:     campaignService,
		analyticsService:    analyticsService,
		contributorService:  contributorService,
		notificationService: notificationService,

		// External dependencies
		paystack:    paystack,
		storage:     storage,
		broadcaster: broadcaster,
		logger:      logger,
		runAsync:    func(f func()) { go f() },
	}
}

func (p *paymentService) InitializeManualPayment(contributorID uint, reference, userEmail, key string) (*models.Payment, error) {

	// validate contributor
	contributor, err := p.contributorService.GetContributorByID(contributorID)
	if err != nil {
		return nil, err
	}

	if contributor.HasPaid() {
		return nil, errs.BadRequest("Contributor has already paid", nil)
	}
	// create a new manual payment
	payment := models.NewManualPayment(contributor.ID, contributor.CampaignID, contributor.GetAmountTotal(), nil)

	// validate user
	//TODO: should campaign creator also provide payment reference 🤔
	if contributor.Email != userEmail {

		if campaign, err := p.campaignService.GetCampaignByID(contributor.CampaignID, key); err != nil {
			return nil, err
		} else if campaign.CreatedBy.Email != userEmail {
			return nil, errs.BadRequest("You are not authorized to perform this action", nil)
		} else {
			payment.SetPaymentStatusToSuccess()
		}
	} else {
		if reference == "" {
			return nil, errs.BadRequest("Reference is required", nil)
		}
	}

	// Upload and update reference to payment proof
	url, id, err := p.storage.UploadFile(reference, "payment/reference")
	if err != nil {
		return nil, errs.InternalServerError(err).Log(p.logger)
	}

	payment.UpdateManualPaymentProof(&models.ManualPaymentProof{
		DocumentURL: url,
		DocumentID:  id,
	})

	// Save Payment
	err = p.repo.Create(payment)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(p.logger)
	}
	contributor.Payment = payment
	// Broadcast event
	p.runAsync(func() {
		p.broadcaster.NewEvent(contributor.CampaignID, websocket.EventTypeContributorUpdated, contributor)
	})

	return payment, nil

}

func (p *paymentService) VerifyPayment(reference string) error {
	// Get the payment
	payment, err := p.repo.GetByReference(reference)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.NotFound("Payment not found")
		}
		return errs.InternalServerError(err).Log(p.logger)
	}

	// Check if the payment has already been verified
	// check if its valid payment method
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

		// Update the contributor
		contributor := payment.Contributor
		contributor.Payment = payment
		go p.broadcaster.NewEvent(contributor.CampaignID, websocket.EventTypeContributorUpdated, contributor)

		go p.notificationService.NotifyPaymentReceived(&contributor, &payment.Campaign)
		return nil

	}

	//TODO: Send email to contributor
	//TODO: set payment status to failed

	return errs.New(fmt.Sprintf("Payment Verification failed :%v ", res.Data.GatewayResponse), http.StatusUnprocessableEntity)

}

// VerifyManualPayment implements interfaces.PaymentService.
func (p *paymentService) VerifyManualPayment(reference, userHandle, key string) error {

	// Validate reference
	payment, err := p.repo.GetByReference(reference)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.NotFound("Payment not found")
		}
		return errs.InternalServerError(err).Log(p.logger)
	}

	// Validate user and campaign creator
	campaign, err := p.campaignService.GetCampaignByID(payment.CampaignID, key)
	if err != nil {
		return err
	}
	if campaign.CreatedBy.Handle != userHandle {
		return errs.BadRequest("Unauthorized: Only campaign creator can verify manual payments", nil)
	}

	// Update payment status
	payment.SetPaymentStatusToSuccess()
	err = p.repo.Update(payment)
	if err != nil {
		return errs.InternalServerError(err).Log(p.logger)
	}
	// Update contributor and broadcast event
	contributor := payment.Contributor
	contributor.Payment = payment

	p.runAsync(func() {
		p.broadcaster.NewEvent(contributor.CampaignID, websocket.EventTypeContributorUpdated, contributor)
	})
	p.runAsync(func() {
		p.notificationService.NotifyPaymentReceived(&contributor, &payment.Campaign)
	})
	p.runAsync(func() {
		p.analyticsService.GetCurrentData().UpdatePaymentStats(payment.PaymentMethod, string(*campaign.FiatCurrency), payment.Amount)
	})

	return nil

}

// InitializePayment implements interfaces.PaymentService.
func (p *paymentService) InitializePayment(contributorID uint, key string) (*models.Payment, error) {

	// Validate the contributor
	contributor, err := p.contributorService.GetContributorByID(contributorID)
	if err != nil {
		return nil, err
	}
	if contributor.HasPaid() {
		return nil, errs.BadRequest("Contributor has already paid", nil)
	}

	// validate campaign
	campaign, err := p.campaignService.GetCampaignByID(contributor.CampaignID, key)

	if err != nil {
		return nil, err
	}
	if campaign.HasEnded() {
		return nil, errs.BadRequest("Campaign has ended", nil)
	}

	// validate payment method
	switch campaign.PaymentMethod {

	case models.PaymentMethodCrypto:
		return nil, errs.BadRequest("Cryptocurrency payment method is not available yet", nil)

	case models.PaymentMethodManual:
		return nil, errs.BadRequest("Campaign payment method is manual", nil)

	case models.PaymentMethodFiat:
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

// TOOD: add notification
// TODO: Cannot test on localHost
// HandlePaystackWebhook implements interfaces.PaymentService.
func (p *paymentService) ProcessPaystackWebhook(event paystack.PaystackWebhookEvent) {
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
		payment.SetPaymentStatusToSuccess()
		err = p.contributorService.UpdateContributor(&contributor)
		if err != nil {
			errs.InternalServerError(err).Log(p.logger)
		}
		return
		// Handle the charge failed event
	case paystack.EventChargeFailed:
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
