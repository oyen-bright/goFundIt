package models

import "time"

type PlatformAnalytics struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Campaign Stats
	TotalCampaigns    int64   `gorm:"not null;default:0" json:"totalCampaigns"`
	NewCampaigns      int64   `gorm:"not null;default:0" json:"newCampaigns"`
	TotalAmountRaised float64 `gorm:"not null;default:0" json:"totalAmountRaised"`

	// User Stats
	TotalUsers  int64 `gorm:"not null;default:0" json:"totalUsers"`
	NewUsers    int64 `gorm:"not null;default:0" json:"newUsers"`
	ActiveUsers int64 `gorm:"not null;default:0" json:"activeUsers"`

	// Activity Stats
	TotalActivities int64 `gorm:"not null;default:0" json:"totalActivities"`
	NewActivities   int64 `gorm:"not null;default:0" json:"newActivities"`

	// Payment Stats
	PaymentMethodStats map[string]int64 `gorm:"serializer:json" json:"paymentMethodStats"`
	FiatCurrencyStats  map[string]int64 `gorm:"serializer:json" json:"fiatCurrencyStats"`
	CryptoTokenStats   map[string]int64 `gorm:"serializer:json" json:"cryptoTokenStats"`

	UpdatedAt time.Time `json:"updatedAt"`
}
