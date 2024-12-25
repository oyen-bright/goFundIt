package interfaces

type AnalyticsService interface {
	StartAnalytics() error
	StopAnalytics()
	ProcessAnalyticsNow() error
}
