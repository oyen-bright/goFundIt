package models

import (
	"fmt"
	"time"

	"github.com/oyen-bright/goFundIt/pkg/utils"
)

type PaymentStatus string

// Payment status constants
const (
	PaymentStatusPending         PaymentStatus = "pending"
	PaymentStatusSucceeded       PaymentStatus = "succeeded"
	PaymentStatusFailed          PaymentStatus = "failed"
	PaymentStatusPendingApproval PaymentStatus = "pending_approval"
)

type PaymentMethod string

const (
	PaymentMethodFiat   PaymentMethod = "fiat"
	PaymentMethodCrypto PaymentMethod = "crypto"
	PaymentMethodManual PaymentMethod = "manual"
)

type FiatCurrency string

const (
	GHS FiatCurrency = "GHS"
	NGN FiatCurrency = "NGN"
)

type CryptoToken string

const (
	USDT CryptoToken = "USDT"
	USDC CryptoToken = "USDC"
	BUSD CryptoToken = "BUSD"
	DAI  CryptoToken = "DAI"
)

// Payment status constants
type Payment struct {
	Reference string `gorm:"type:text;primaryKey" json:"reference"`

	ContributorID uint   `gorm:"not null;foreignKey:ContributorID" json:"contributorId"`
	CampaignID    string `gorm:"not null;foreignKey:CampaignID" json:"campaignId"`

	Amount          float64             `gorm:"not null;type:numeric(10,2)" json:"amount"`
	PaymentMethod   PaymentMethod       `gorm:"not null;size:50" json:"paymentMethod"`
	PaymentStatus   PaymentStatus       `gorm:"not null;size:50;default:'pending'" json:"paymentStatus"`
	GatewayResponse *string             `gorm:"type:jsonb" json:"gatewayResponse,omitempty"`
	CreatedAt       time.Time           `gorm:"default:CURRENT_TIMESTAMP;index" json:"createdAt"`
	UpdatedAt       time.Time           `gorm:"default:CURRENT_TIMESTAMP;index" json:"-"`
	PaymentProof    *ManualPaymentProof `gorm:"embedded" json:"paymentProof,omitempty"`

	// Relations
	Contributor Contributor `gorm:"foreignKey:ContributorID;references:ID" json:"-"`
	Campaign    Campaign    `gorm:"foreignKey:CampaignID;references:ID" json:"-"`

	//PaymentURL
	AuthorizationURL string `gorm:"-" json:"authorization_url"`
}

// ManualPaymentProof represents proof of payment for manual payments
type ManualPaymentProof struct {
	DocumentID  string `json:"-"`
	DocumentURL string `json:"url"`
}

// Constructor

// NewPayment creates a new Payment instance with the provided parameters
func NewPayment(contributorID uint, campaignID, reference string, amount float64, paymentMethod PaymentMethod, GatewayResponse string, authorizationURL string) *Payment {
	return &Payment{
		ContributorID:    contributorID,
		CampaignID:       campaignID,
		Reference:        reference,
		Amount:           amount,
		PaymentMethod:    paymentMethod,
		PaymentStatus:    PaymentStatusPending,
		GatewayResponse:  &GatewayResponse,
		AuthorizationURL: authorizationURL,
	}
}

func NewManualPayment(contributorID uint, campaignID string, amount float64, paymentProof *ManualPaymentProof) *Payment {
	return &Payment{
		ContributorID: contributorID,
		CampaignID:    campaignID,
		Reference:     generateManualReference(contributorID),
		Amount:        amount,
		PaymentProof:  paymentProof,
		PaymentMethod: PaymentMethodManual,
		PaymentStatus: PaymentStatusPendingApproval,
	}
}

func NewFiatPayment(contributorID uint, campaignID, reference string, amount float64, authorizationURL string) *Payment {

	//TODO: add charges so that the amount is not the same as the contributor's total amount
	return &Payment{
		ContributorID:    contributorID,
		CampaignID:       campaignID,
		Reference:        reference,
		Amount:           amount,
		PaymentMethod:    PaymentMethodFiat,
		PaymentStatus:    PaymentStatusPending,
		AuthorizationURL: authorizationURL,
	}
}

// SetPaymentStatusToFailed updates the payment status to failed
func (p *Payment) SetPaymentStatusToFailed() {
	p.PaymentStatus = PaymentStatusFailed
}

// SetPaymentStatusToSuccess updates the payment status to succeeded
func (p *Payment) SetPaymentStatusToSuccess() {
	p.PaymentStatus = PaymentStatusSucceeded
}

// GetPaymentLink returns the payment link for the payment
func (p *Payment) GetPaymentLink() interface{} {
	return map[string]interface{}{
		"reference":     p.Reference,
		"paymentLink":   p.AuthorizationURL,
		"paymentStatus": p.PaymentStatus,
	}
}

// UpdateManualPaymentProof updates the payment proof for manual payments
func (p *Payment) UpdateManualPaymentProof(proof *ManualPaymentProof) {
	p.PaymentProof = proof

}

// Helper Functions --------------------------------------------------------------------

func generateManualReference(contributorID uint) string {
	return utils.GenerateRandomAlphaNumeric(fmt.Sprintf("M-%d", contributorID), 8)
}
