package postgress

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func (r *analyticsRepository) Save(analytics *models.PlatformAnalytics) error {
	return r.db.Save(analytics).Error
}

func (r *analyticsRepository) Get(date time.Time) (*models.PlatformAnalytics, error) {
	var analytics models.PlatformAnalytics

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	err := r.db.Where("DATE(created_at) = ?", startOfDay.Format("2006-01-02")).First(&analytics).Error

	if err == gorm.ErrRecordNotFound {
		analytics = models.PlatformAnalytics{
			FiatStats:      make(map[string]models.CurrencyStats),
			CryptoStats:    make(map[string]models.CurrencyStats),
			PaymentMethods: models.PaymentMethodStats{},
			CreatedAt:      date,
			UpdatedAt:      date,
		}

		if err := r.db.Create(&analytics).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &analytics, nil
}

func NewAnalyticsRepository(db *gorm.DB) *analyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}
