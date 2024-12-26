// TODO: refactor the whole service
package services

import (
	"fmt"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	"github.com/oyen-bright/goFundIt/pkg/errs"

	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/robfig/cron/v3"
)

type analyticsService struct {
	repo       repositories.AnalyticsRepository
	emailer    email.Emailer
	adminEmail string
	logger     logger.Logger
	data       *models.PlatformAnalytics
	cron       *cron.Cron
}

// AnalyticsService interface implementation
func NewAnalyticsService(
	analyticsRepo repositories.AnalyticsRepository,
	analyticsReportEmail string,
	emailer email.Emailer,
	logger logger.Logger,

) services.AnalyticsService {
	service := &analyticsService{
		adminEmail: analyticsReportEmail,
		repo:       analyticsRepo,
		emailer:    emailer,
		logger:     logger,
	}

	service.data = service.getCurrentData()
	return service
}

// StartAnalytics starts the cron job for daily analytics at 23:00 UTC
func (s *analyticsService) StartAnalytics() error {
	s.cron = cron.New()

	// Schedule analytics processing for 23:00 UTC daily
	_, err := s.cron.AddFunc("0 23 * * *", func() {
		s.processDailyAnalytics()
	})

	if err != nil {
		return fmt.Errorf("failed to setup analytics cron: %w", err)
	}

	s.cron.Start()
	s.logger.Info("Analytics service started - scheduled for 23:00 UTC daily", nil)
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
	err := s.processDailyAnalytics()

	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

// GetCurrentData implements interfaces.AnalyticsService.
func (s *analyticsService) GetCurrentData() *models.PlatformAnalytics {
	return s.data
}

func (s *analyticsService) processDailyAnalytics() error {
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)

	err := s.repo.Save(s.data)
	if err != nil {
		return fmt.Errorf("failed to get today's analytics: %w", err)
	}

	yesterdayStats, err := s.repo.Get(yesterday)
	if err != nil {
		return fmt.Errorf("failed to get yesterday's analytics: %w", err)
	}

	// Send daily report
	if err := s.sendDailyReport(s.data, yesterdayStats, now); err != nil {
		return fmt.Errorf("failed to send daily report: %w", err)
	}

	s.logger.Info("Daily analytics processed successfully", map[string]interface{}{
		"date": now.Format("2006-01-02"),
	})

	return nil
}

func (s *analyticsService) sendDailyReport(
	today *models.PlatformAnalytics,
	yesterday *models.PlatformAnalytics,
	reportDate time.Time,
) error {
	template := emailTemplates.AnalyticsReport(
		[]string{s.adminEmail},
		today,
		today.GenerateComparison(yesterday),
		reportDate,
	)

	if err := s.emailer.SendEmailTemplate(*template); err != nil {
		return fmt.Errorf("failed to send analytics report: %w", err)
	}

	s.logger.Info("Daily analytics report sent", map[string]interface{}{
		"date":  reportDate.Format("2006-01-02"),
		"email": s.adminEmail,
	})

	return nil
}

// Helper function to get the current analytics data
func (s *analyticsService) getCurrentData() *models.PlatformAnalytics {
	data, err := s.repo.Get(time.Now().UTC())
	if err != nil {
		s.logger.Error(err, "Failed to get current analytics data", nil)
		return &models.PlatformAnalytics{}
	}

	return data
}
