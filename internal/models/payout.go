package models

import (
	"errors"
	"time"

	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

type PayoutStatus string

type FiatAccount struct {
	BankCode      string `gorm:"size:10" json:"bankCode"`
	BankName      string `gorm:"size:100" json:"bankName"`
	AccountName   string `gorm:"size:100" json:"accountName"`
	AccountNumber string `gorm:"size:20" json:"accountNumber"`
	Currency      string `gorm:"size:10" json:"currency"`
}

type CryptoAccount struct {
	CryptoToken CryptoToken `gorm:"size:10" json:"cryptoToken"`
	Address     string      `gorm:"size:255" json:"address"`
}

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

type Payout struct {
	ID            string         `gorm:"primaryKey;size:255" json:"-"`
	RecipientID   string         `gorm:"size:255" json:"-"`
	CampaignID    string         `gorm:"not null;foreignKey:CampaignID" json:"campaignId"`
	Amount        float64        `gorm:"not null;type:numeric(10,2)" json:"amount"`
	PayoutMethod  PaymentMethod  `gorm:"not null;size:50" json:"payoutMethod"`
	Status        PayoutStatus   `gorm:"not null;size:50;default:'pending'" json:"status"`
	Reference     string         `gorm:"size:255" json:"reference"`
	FiatAccount   *FiatAccount   `gorm:"embedded" json:"fiatAccount,omitempty"`
	CryptoAccount *CryptoAccount `gorm:"embedded" json:"cryptoAccount,omitempty"`
	FailureReason *string        `gorm:"size:255" json:"failureReason"`
	ProcessedAt   *string        `gorm:"size:255" json:"processedAt"`
	CompletedAt   *string        `gorm:"size:255" json:"completedAt"`
	CreatedAt     string         `gorm:"default:CURRENT_TIMESTAMP;index" json:"-"`
	UpdatedAt     string         `gorm:"default:CURRENT_TIMESTAMP;index" json:"-"`
}

// NewPayout creates a new payout instance
func NewPayout(campaignID string, amount float64, payoutMethod PaymentMethod) *Payout {
	return &Payout{
		ID:           generatePayoutId(),
		CampaignID:   campaignID,
		Amount:       amount,
		PayoutMethod: payoutMethod,
	}
}

// NewFiatPayout creates a new fiat payout instance
func NewFiatPayout(campaignID string, amount float64, bankCode, bankName, accountName, accountNumber, currency, recipientId string) *Payout {
	return &Payout{
		ID:           generatePayoutId(),
		CampaignID:   campaignID,
		Amount:       amount,
		RecipientID:  recipientId,
		PayoutMethod: PaymentMethodFiat,

		FiatAccount: &FiatAccount{
			Currency:      currency,
			BankCode:      bankCode,
			BankName:      bankName,
			AccountName:   accountName,
			AccountNumber: accountNumber,
		},
	}
}

// NewManualPayout creates a new manual payout instance
func NewManualPayout(campaignID string, amount float64, recipientId string) *Payout {
	return &Payout{
		ID:           generatePayoutId(),
		CampaignID:   campaignID,
		Amount:       amount,
		RecipientID:  recipientId,
		PayoutMethod: PaymentMethodManual,
	}
}

// BeforeCreate ensures ID is not null before creating a new payout
func (p *Payout) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		return errors.New("ID is required")
	}
	return nil
}

// MarkPayoutCompleted sets the payout status as completed and updates the completion time
func (p *Payout) MarkPayoutCompleted() {
	p.Status = PayoutStatusCompleted
	now := time.Now().UTC().Format(time.RFC3339)
	p.CompletedAt = &now
}

// MarkPayoutFailed sets the payout status as failed and updates the failure reason
func (p *Payout) MarkPayoutFailed(reason string) {
	p.Status = PayoutStatusFailed
	p.FailureReason = &reason
}

// MarkPayoutProcessing sets the payout status as processing
func (p *Payout) MarkPayoutProcessing() {
	p.Status = PayoutStatusProcessing
	now := time.Now().UTC().Format(time.RFC3339)
	p.ProcessedAt = &now
}

// NewCryptoPayout creates a new crypto payout instance
func NewCryptoPayout(campaignID string, amount float64, cryptoToken CryptoToken, address string) *Payout {
	return &Payout{
		CampaignID:   campaignID,
		Amount:       amount,
		PayoutMethod: PaymentMethodCrypto,
		CryptoAccount: &CryptoAccount{
			CryptoToken: cryptoToken,
			Address:     address,
		},
	}
}

// Helper Functions
func generatePayoutId() string {
	return utils.GenerateRandomAlphaNumeric("PYT-", 12)
}
