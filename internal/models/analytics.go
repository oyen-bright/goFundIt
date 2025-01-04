package models

import "time"

type PaymentMethodStats struct {
	Manual int64 `json:"manual"`
	Crypto int64 `json:"crypto"`
	Fiat   int64 `json:"fiat"`
}

type CurrencyStats struct {
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

type AnalyticsComparison struct {
	Users struct {
		Total      int64   `json:"total"`
		Change     int64   `json:"change"`
		Percentage float64 `json:"percentage"`
	} `json:"users"`
	Campaigns struct {
		Total      int64   `json:"total"`
		Change     int64   `json:"change"`
		Percentage float64 `json:"percentage"`
	} `json:"campaigns"`
	Activities struct {
		Total      int64   `json:"total"`
		Change     int64   `json:"change"`
		Percentage float64 `json:"percentage"`
	} `json:"activities"`
	Finances struct {
		TotalRaised float64 `json:"total_raised"`
		Change      float64 `json:"change"`
		Percentage  float64 `json:"percentage"`
	} `json:"finances"`
}

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
	PaymentMethods PaymentMethodStats       `gorm:"serializer:json" json:"paymentMethodStats"`
	FiatStats      map[string]CurrencyStats `gorm:"serializer:json" json:"fiatCurrencyStats"`
	CryptoStats    map[string]CurrencyStats `gorm:"serializer:json" json:"cryptoTokenStats"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IncrementCampaigns increases campaign-related counters
func (pa *PlatformAnalytics) IncrementCampaigns(amount float64) {
	pa.TotalCampaigns++
	pa.NewCampaigns++
	pa.TotalAmountRaised += amount
	pa.UpdatedAt = time.Now().UTC()
}

// IncrementUsers increases user-related counters
func (pa *PlatformAnalytics) IncrementUsers(isActive bool) {
	pa.TotalUsers++
	pa.NewUsers++
	if isActive {
		pa.ActiveUsers++
	}
	pa.UpdatedAt = time.Now().UTC()
}

// IncrementUsers increases user-related counters
func (pa *PlatformAnalytics) IncrementMultipleUsers(count int64, isActive bool) {
	pa.TotalUsers += count
	pa.NewUsers += count
	if isActive {
		pa.ActiveUsers += count
	}
	pa.UpdatedAt = time.Now().UTC()
}

// IncrementActivities increases activity counters
func (pa *PlatformAnalytics) IncrementActivities() {
	pa.TotalActivities++
	pa.NewActivities++
	pa.UpdatedAt = time.Now().UTC()
}

// UpdatePaymentStats updates payment-related statistics
func (pa *PlatformAnalytics) UpdatePaymentStats(paymentType PaymentMethod, currency string, amount float64) {
	// Initialize maps if nil
	if pa.FiatStats == nil {
		pa.FiatStats = make(map[string]CurrencyStats)
	}
	if pa.CryptoStats == nil {
		pa.CryptoStats = make(map[string]CurrencyStats)
	}

	switch paymentType {
	case PaymentMethodManual:
		pa.PaymentMethods.Manual++
	case PaymentMethodCrypto:
		pa.PaymentMethods.Crypto++
		stats := pa.CryptoStats[currency]
		stats.Amount += amount
		stats.Count++
		pa.CryptoStats[currency] = stats
	case PaymentMethodFiat:
		pa.PaymentMethods.Fiat++
		stats := pa.FiatStats[currency]
		stats.Amount += amount
		stats.Count++
		pa.FiatStats[currency] = stats
	}

	pa.UpdatedAt = time.Now().UTC()
}

func (pa *PlatformAnalytics) ResetNewStats() {
	pa.NewCampaigns = 0
	pa.NewUsers = 0
	pa.NewActivities = 0
	pa.UpdatedAt = time.Now().UTC()
}

func (pa *PlatformAnalytics) GenerateComparison(yesterday *PlatformAnalytics) AnalyticsComparison {
	var comparison AnalyticsComparison

	// Users comparison
	comparison.Users.Total = pa.TotalUsers
	comparison.Users.Change = pa.NewUsers - yesterday.NewUsers
	if yesterday.NewUsers > 0 {
		comparison.Users.Percentage = float64(pa.NewUsers-yesterday.NewUsers) / float64(yesterday.NewUsers) * 100
	}

	// Campaigns comparison
	comparison.Campaigns.Total = pa.TotalCampaigns
	comparison.Campaigns.Change = pa.NewCampaigns - yesterday.NewCampaigns
	if yesterday.NewCampaigns > 0 {
		comparison.Campaigns.Percentage = float64(pa.NewCampaigns-yesterday.NewCampaigns) / float64(yesterday.NewCampaigns) * 100
	}

	// Activities comparison
	comparison.Activities.Total = pa.TotalActivities
	comparison.Activities.Change = pa.NewActivities - yesterday.NewActivities
	if yesterday.NewActivities > 0 {
		comparison.Activities.Percentage = float64(pa.NewActivities-yesterday.NewActivities) / float64(yesterday.NewActivities) * 100
	}

	// Finances comparison
	comparison.Finances.TotalRaised = pa.TotalAmountRaised
	todayRaised := pa.TotalAmountRaised - yesterday.TotalAmountRaised
	comparison.Finances.Change = todayRaised
	if yesterday.TotalAmountRaised > 0 {
		comparison.Finances.Percentage = (todayRaised / yesterday.TotalAmountRaised) * 100
	}

	return comparison
}
