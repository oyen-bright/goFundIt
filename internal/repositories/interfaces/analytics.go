package interfaces

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
)

type AnalyticsRepository interface {
	Save(analytics *models.PlatformAnalytics) error
	Get(date time.Time) (*models.PlatformAnalytics, error)
}
