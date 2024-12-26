// TODO: refactor
package services

import (
	"fmt"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"

	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/robfig/cron/v3"
)

type analyticsService struct {
	campaignService services.CampaignService
	authService     services.AuthService
	analyticsRepo   repositories.AnalyticsRepository
	emailer         email.Emailer
	adminEmail      string
	logger          logger.Logger
	cron            *cron.Cron
}

// AnalyticsService interface implementation
func NewAnalyticsService(
	campaignService services.CampaignService,
	authService services.AuthService,
	analyticsRepo repositories.AnalyticsRepository,
	analyticsReportEmail string,
	emailer email.Emailer,
	logger logger.Logger,

) services.AnalyticsService {
	return &analyticsService{
		campaignService: campaignService,
		authService:     authService,
		adminEmail:      analyticsReportEmail,
		analyticsRepo:   analyticsRepo,
		emailer:         emailer,
		logger:          logger,
	}
}

// StartAnalytics starts the cron job for daily analytics at 23:00 UTC
func (s *analyticsService) StartAnalytics() error {
	s.cron = cron.New()

	// Schedule analytics processing for 23:00 UTC daily
	_, err := s.cron.AddFunc("0 23 * * *", func() {
		if err := s.processDailyAnalytics(); err != nil {
			s.logger.Error(err, "error processDailyAnalytics ", nil)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to setup analytics cron: %w", err)
	}

	s.cron.Start()
	return nil
}

// StopAnalytics stops the cron job
func (s *analyticsService) StopAnalytics() {
	if s.cron != nil {
		s.cron.Stop()
		s.logger.Info("Analytics service stopped", nil)
	}
}

// ProcessAnalyticsNow triggers an immediate analytics processing
func (s *analyticsService) ProcessAnalyticsNow() error {
	return s.processDailyAnalytics()
}

// Private methods

func (s *analyticsService) processDailyAnalytics() error {
	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)

	// Get data for analysis
	campaigns, err := s.campaignService.GetCampaignsForAnalytics(yesterday, now)
	if err != nil {
		return fmt.Errorf("failed to get campaigns: %w", err)
	}

	users, err := s.authService.GetUsersByCreatedDateRange(yesterday, now)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// Process today's data
	todayAnalytics := s.createDailyAnalytics(campaigns, users, yesterday)

	// Get and update current analytics
	currentAnalytics, err := s.analyticsRepo.Get()
	if err != nil {
		return fmt.Errorf("failed to get current analytics: %w", err)
	}

	// Generate comparison report
	comparison := s.generateComparisonReport(currentAnalytics, todayAnalytics)

	// Update and save analytics
	updatedAnalytics := s.updateAnalytics(currentAnalytics, todayAnalytics)
	if err := s.analyticsRepo.Save(updatedAnalytics); err != nil {
		return fmt.Errorf("failed to save analytics: %w", err)
	}

	// Send daily report
	if err := s.sendDailyReport(todayAnalytics, comparison, now); err != nil {
		return fmt.Errorf("failed to send daily report: %w", err)
	}

	return nil
}

func (s *analyticsService) createDailyAnalytics(campaigns []models.Campaign, users []models.User, yesterday time.Time) *models.PlatformAnalytics {
	analytics := &models.PlatformAnalytics{
		PaymentMethodStats: make(map[string]int64),
		FiatCurrencyStats:  make(map[string]int64),
		CryptoTokenStats:   make(map[string]int64),
	}

	// Process campaigns
	for _, campaign := range campaigns {
		// Campaign counts
		analytics.TotalCampaigns++
		if campaign.CreatedAt.After(yesterday) {
			analytics.NewCampaigns++
		}

		// Financial stats
		campaignRaised := float64(0)
		for _, contributor := range campaign.Contributors {
			campaignRaised += contributor.Amount
		}
		analytics.TotalAmountRaised += campaignRaised

		// Activity stats
		for _, activity := range campaign.Activities {
			analytics.TotalActivities++
			if activity.CreatedAt.After(yesterday) {
				analytics.NewActivities++
			}
		}

		// Payment method stats
		analytics.PaymentMethodStats[string(campaign.PaymentMethod)]++
		if campaign.FiatCurrency != nil {
			analytics.FiatCurrencyStats[string(*campaign.FiatCurrency)]++
		}
		if campaign.CryptoToken != nil {
			analytics.CryptoTokenStats[string(*campaign.CryptoToken)]++
		}
	}

	// Process users
	analytics.TotalUsers = int64(len(users))
	for _, user := range users {
		if user.CreatedAt.After(yesterday) {
			analytics.NewUsers++
		}
		if user.UpdatedAt.After(yesterday) {
			analytics.ActiveUsers++
		}
	}

	return analytics
}

func (s *analyticsService) updateAnalytics(current, today *models.PlatformAnalytics) *models.PlatformAnalytics {
	if current == nil {
		return today
	}

	// Update base stats
	current.TotalCampaigns += today.NewCampaigns
	current.TotalAmountRaised += today.TotalAmountRaised
	current.TotalActivities += today.NewActivities
	current.TotalUsers += today.NewUsers

	// Update maps
	for method, count := range today.PaymentMethodStats {
		current.PaymentMethodStats[method] += count
	}
	for currency, count := range today.FiatCurrencyStats {
		current.FiatCurrencyStats[currency] += count
	}
	for token, count := range today.CryptoTokenStats {
		current.CryptoTokenStats[token] += count
	}

	// Set today's new counts
	current.NewCampaigns = today.NewCampaigns
	current.NewActivities = today.NewActivities
	current.NewUsers = today.NewUsers
	current.ActiveUsers = today.ActiveUsers

	return current
}

func (s *analyticsService) generateComparisonReport(current, today *models.PlatformAnalytics) map[string]interface{} {
	return map[string]interface{}{
		"campaigns": map[string]interface{}{
			"total":     current.TotalCampaigns,
			"new_today": today.NewCampaigns,
		},
		"users": map[string]interface{}{
			"total":        current.TotalUsers,
			"new_today":    today.NewUsers,
			"active_today": today.ActiveUsers,
		},
		"activities": map[string]interface{}{
			"total":     current.TotalActivities,
			"new_today": today.NewActivities,
		},
		"finances": map[string]interface{}{
			"total_raised": current.TotalAmountRaised,
			"raised_today": today.TotalAmountRaised,
		},
	}
}

func (s *analyticsService) sendDailyReport(
	today *models.PlatformAnalytics,
	comparison map[string]interface{},
	reportDate time.Time,
) error {

	template := emailTemplates.AnalyticsReport([]string{s.adminEmail}, today, comparison, reportDate)

	return s.emailer.SendEmailTemplate(*template)
}
