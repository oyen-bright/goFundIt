package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) interfaces.AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}

// Save updates or creates the analytics record
func (r *analyticsRepository) Save(analytics *models.PlatformAnalytics) error {
	// Initialize maps if they're nil
	if analytics.PaymentMethodStats == nil {
		analytics.PaymentMethodStats = make(map[string]int64)
	}
	if analytics.FiatCurrencyStats == nil {
		analytics.FiatCurrencyStats = make(map[string]int64)
	}
	if analytics.CryptoTokenStats == nil {
		analytics.CryptoTokenStats = make(map[string]int64)
	}

	// Get existing record or create new one
	var existing models.PlatformAnalytics
	result := r.db.First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new record if none exists
		return r.db.Create(analytics).Error
	}

	// Update existing record
	analytics.ID = existing.ID
	return r.db.Save(analytics).Error
}

// Get retrieves the analytics record
func (r *analyticsRepository) Get() (*models.PlatformAnalytics, error) {
	var analytics models.PlatformAnalytics

	err := r.db.First(&analytics).Error
	if err == gorm.ErrRecordNotFound {
		return &models.PlatformAnalytics{
			PaymentMethodStats: make(map[string]int64),
			FiatCurrencyStats:  make(map[string]int64),
			CryptoTokenStats:   make(map[string]int64),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &analytics, nil
}
