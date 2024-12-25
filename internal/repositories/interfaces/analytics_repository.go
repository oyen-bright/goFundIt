package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type AnalyticsRepository interface {
	Save(analytics *models.PlatformAnalytics) error
	Get() (*models.PlatformAnalytics, error)
}
