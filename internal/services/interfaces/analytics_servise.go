package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type AnalyticsService interface {
	StartAnalytics() error
	StopAnalytics()
	ProcessAnalyticsNow() error
	GetCurrentData() *models.PlatformAnalytics
}
