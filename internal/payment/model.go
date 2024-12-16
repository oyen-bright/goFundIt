package payment

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/campaign"
	"github.com/oyen-bright/goFundIt/internal/contributor"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID            uint                    `gorm:"primaryKey"`
	ContributorID uint                    `gorm:"not null"`
	CampaignID    string                  `gorm:"not null;index"`
	Amount        float64                 `gorm:"not null;type:numeric(10,2)"`
	PaymentMethod string                  `gorm:"not null;size:50"`
	PaymentStatus PaymentStatus           `gorm:"not null;size:50;default:'pending'"`
	TransactionID string                  `gorm:"unique;size:255"`
	CreatedAt     time.Time               `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time               `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Contributor   contributor.Contributor `gorm:"foreignKey:ContributorID;references:ID"`
	Campaign      campaign.Campaign       `gorm:"foreignKey:CampaignID;references:ID"`
}
