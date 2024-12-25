package services

import (
	"fmt"
	"log"

	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

type payoutService struct {
	repo                interfaces.PayoutRepository
	campaignService     services.CampaignService
	notificationService services.NotificationService
	paystack            *paystack.Client
	broadCaster         services.EventBroadcaster
	logger              logger.Logger
}

// NewPayoutService creates a new instance of the payout service
func NewPayoutService(payoutRepo interfaces.PayoutRepository, campaignService services.CampaignService, notificationService services.NotificationService, paystack *paystack.Client, broadCaster services.EventBroadcaster, logger logger.Logger) services.PayoutService {
	return &payoutService{
		broadCaster:         broadCaster,
		campaignService:     campaignService,
		notificationService: notificationService,
		repo:                payoutRepo,
		paystack:            paystack,
		logger:              logger}
}

// InitializeManualPayout implements interfaces.PayoutService.
func (p *payoutService) InitializeManualPayout(campaignID string, userHandle string) (*models.Payout, error) {
	// Validate the campaign and user
	campaign, err := p.campaignService.GetCampaignByIDWithContributors(campaignID)
	if err != nil {
		return nil, err
	}
	if campaign.CreatedBy.Handle != userHandle {
		return nil, errs.BadRequest("You are not authorized to perform this action", nil)
	}

	// Validate payout status
	if campaign.Payout != nil {
		if campaign.Payout.Status == models.PayoutStatusPending {
			return nil, errs.BadRequest("You have a pending payout", nil)
		}

		if campaign.Payout.Status == models.PayoutStatusCompleted {
			return nil, errs.BadRequest("You have already completed a payout", nil)
		}
	}
	if !campaign.CanInitiatePayout() {
		return nil, errs.BadRequest("Cannot initiate payout: Some contributors haven't completed their payments. Please ensure all contributors have paid or remove unpaid contributors before proceeding.", nil)
	}

	// Process Payout
	payout := models.NewManualPayout(campaignID, campaign.GetPayoutAmount(), "")
	payout.MarkPayoutCompleted()

	// Create Payout
	if err := p.repo.Create(payout); err != nil {
		return nil, err
	}

	// Broadcast Payout
	p.broadCaster.NewEvent(campaignID, websocket.EventTypePayoutUpdated, payout)

	go p.notificationService.NotifyPayoutCollected(campaign)

	return payout, nil

}

// InitializePayout implements interfaces.PayoutService.
func (p *payoutService) InitializePayout(campaignID string, userHandle string, req dto.PayoutRequest) (*models.Payout, error) {
	// Validate the campaign and user
	campaign, err := p.campaignService.GetCampaignByIDWithContributors(campaignID)
	if err != nil {
		return nil, err
	}
	if campaign.CreatedBy.Handle != userHandle {
		return nil, errs.BadRequest("You are not authorized to perform this action", nil)
	}

	// Validate payout status
	if campaign.Payout != nil {
		if campaign.Payout.Status == models.PayoutStatusPending {
			return nil, errs.BadRequest("You have a pending payout", nil)
		}

		if campaign.Payout.Status == models.PayoutStatusCompleted {
			return nil, errs.BadRequest("You have already completed a payout", nil)
		}
	}
	if !campaign.CanInitiatePayout() {
		return nil, errs.BadRequest("Cannot initiate payout: Some contributors haven't completed their payments. Please ensure all contributors have paid or remove unpaid contributors before proceeding.", nil)
	}

	// Process Payout
	var payout models.Payout
	switch campaign.PaymentMethod {

	case models.PaymentMethodCrypto:
		return nil, errs.BadRequest("Cryptocurrency payout  is not available yet", nil)

	case models.PaymentMethodManual:
		return nil, errs.BadRequest("Campaign payment method is manual", nil)

	case models.PaymentMethodFiat:
		// Create Recipient
		transferRecipient := paystack.NewRecipient(req.AccountName, req.AccountNumber, req.BankCode, string(*campaign.FiatCurrency))
		res, err := p.paystack.CreateRecipient(*transferRecipient)

		if err != nil {
			return nil, errs.InternalServerError(err).Log(p.logger)
		}
		if !res.Status {
			return nil, errs.InternalServerError(err).Log(p.logger)
		}
		payout.MarkPayoutProcessing()
		payout = *models.NewFiatPayout(campaignID, campaign.GetPayoutAmount(), req.BankCode, req.BankName, req.AccountName, req.AccountNumber, string(*campaign.FiatCurrency), res.Data.RecipientCode)

	}

	// Create Payout
	if err := p.repo.Create(&payout); err != nil {
		return nil, err
	}

	//Process Transfer
	go p.processPayoutTransfer(payout)
	return &payout, nil
}

// VerifyAccount implements interfaces.PayoutService.
func (p *payoutService) VerifyAccount(req dto.VerifyAccountRequest) (interface{}, error) {
	res, err := p.paystack.ResolveAccount(req.AccountNumber, req.BankCode)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(p.logger)
	}

	if !res.Status {
		return nil, errs.BadRequest("Account verification failed", nil)
	}
	return res.Data, nil

}

// GetBankList implements interfaces.PayoutService.
func (p *payoutService) GetBankList() ([]interface{}, error) {
	bankListResponse, err := p.paystack.GetBanks()
	log.Println(bankListResponse)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(p.logger)
	}
	if !bankListResponse.Status {
		return nil, errs.InternalServerError(errs.New("Failed to get bank list")).Log(p.logger)
	}

	banks := make([]interface{}, len(bankListResponse.Data))
	for i, bank := range bankListResponse.Data {
		banks[i] = bank
	}
	return banks, nil

}

// GetPayoutByCampaignID implements interfaces.PayoutService.
func (p *payoutService) GetPayoutByCampaignID(campaignID string) (*models.Payout, error) {
	payout, _, err := p.repo.GetByCampaignID(campaignID, 1, 0)
	if err != nil {

		if database.Error(err).IsNotfound() {
			return nil, errs.NotFound("No payout found for this campaign")
		}
		return nil, errs.InternalServerError(err).Log(p.logger)
	}
	return &payout[0], nil
}

// Helper function to process payout transfer

// ProcessPayoutTransfer
func (p *payoutService) processPayoutTransfer(payout models.Payout) {
	switch payout.PayoutMethod {
	case models.PaymentMethodFiat:
		p.processFiatTransfer(payout)
		return
	case models.PaymentMethodCrypto:
		return
	case models.PaymentMethodManual:
		return
	}

}

// ProcessFiatTransfer
func (p *payoutService) processFiatTransfer(payout models.Payout) {
	transfer := paystack.NewTransfer(fmt.Sprint("Payout for campaign: ", payout.CampaignID), payout.RecipientID, payout.FiatAccount.Currency, payout.Amount)
	res, err := p.paystack.InitiateTransfer(*transfer)
	if err != nil {
		payout.MarkPayoutFailed(err.Error())
		p.repo.Update(&payout)
		p.broadCaster.NewEvent(payout.CampaignID, websocket.EventTypePayoutUpdated, payout)
		return
	}
	if !res.Status {
		payout.MarkPayoutFailed(res.Message)
		p.repo.Update(&payout)
		p.broadCaster.NewEvent(payout.CampaignID, websocket.EventTypePayoutUpdated, payout)
		return
	}
	payout.MarkPayoutProcessing()
	p.repo.Update(&payout)
	p.broadCaster.NewEvent(payout.CampaignID, websocket.EventTypePayoutUpdated, payout)
}
