package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/robfig/cron/v3"
)

type cronService struct {
	campaignService     interfaces.CampaignService
	notificationService interfaces.NotificationService
	logger              logger.Logger
	cron                *cron.Cron
}

func NewCronService(campaignService interfaces.CampaignService, notificationService interfaces.NotificationService, logger logger.Logger) interfaces.CronService {
	return &cronService{
		campaignService:     campaignService,
		cron:                cron.New(cron.WithLocation(time.UTC)),
		notificationService: notificationService,
		logger:              logger,
	}
}

func (n *cronService) StartCronJobs() error {
	// Daily cleanup at midnight UTC
	_, err := n.cron.AddFunc("0 0 * * *", func() {
		monitorCronJob("cleanup-campaign", func() {
			n.cleanUpExpiredCampaign()
		})
	})
	n.logger.Info("Cleanup job scheduled - running at midnight UTC daily", nil)
	if err != nil {
		return fmt.Errorf("failed to schedule cleanup job: %w", err)
	}

	// Check contribution reminders every 3 days
	_, err = n.cron.AddFunc("0 0 */3 * *", func() {
		monitorCronJob("contribution-reminders", func() {
			n.checkContributionReminders()
		})
	})
	n.logger.Info("Contribution reminders job scheduled - running every 3 days at midnight UTC", nil)
	if err != nil {
		return fmt.Errorf("failed to schedule contribution reminders job: %w", err)
	}

	// Check campaign deadlines every day
	_, err = n.cron.AddFunc("0 0 * * *", func() {
		monitorCronJob("campaign-deadline", func() {
			n.checkCampaignDeadline()
		})
	})
	n.logger.Info("Campaign deadline check job scheduled - running at midnight UTC daily", nil)

	if err != nil {
		return fmt.Errorf("failed to schedule campaign deadline reminders job: %w", err)
	}

	return nil
}

// monitorCronJob wraps a cron job with Sentry monitoring
func monitorCronJob(slug string, job func()) {
	// Start the job
	checkInId := sentry.CaptureCheckIn(
		&sentry.CheckIn{
			MonitorSlug: slug,
			Status:      sentry.CheckInStatusInProgress,
		},
		nil,
	)

	startTime := time.Now()

	// Execute job with panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Log error to Sentry
			sentry.CurrentHub().Recover(r)

			// Mark job as failed
			sentry.CaptureCheckIn(
				&sentry.CheckIn{
					ID:          *checkInId,
					MonitorSlug: slug,
					Status:      sentry.CheckInStatusError,
					Duration:    time.Since(startTime),
				},
				nil,
			)
		}
	}()

	// Run the job
	job()

	// Mark job as complete
	sentry.CaptureCheckIn(
		&sentry.CheckIn{
			ID:          *checkInId,
			MonitorSlug: slug,
			Status:      sentry.CheckInStatusOK,
			Duration:    time.Since(startTime),
		},
		nil,
	)
}

func (n *cronService) StopCronJobs() {
	if n.cron != nil {
		n.cron.Stop()
	}
}

// cleanUpExpiredCampaign checks if a campaign has ended and if it can be cleaned up
func (n *cronService) cleanUpExpiredCampaign() {

	expiredCampaigns, err := n.campaignService.GetExpiredCampaigns()
	if err != nil {
		return
	}

	for _, campaign := range expiredCampaigns {

		go func(campaign *models.Campaign) {
			if !campaign.HasEnded() {
				return
			}

			if !campaign.HasEnded() {
				return
			}

			if !campaign.CanCleanUp() {
				n.notificationService.NotifyCampaignPayoutRequired(campaign)
				return
			}

			campaignFullData, err := n.campaignService.GetCampaignByIDWithAllRelatedData(campaign.ID)
			if err != nil {
				return
			}
			filePath, err := createJSONExport(*campaignFullData)
			if err != nil {
				return
			}
			defer os.Remove(filePath)
			err = n.notificationService.NotifyCampaignCleanUp(campaign, filePath)
			if err != nil {
				return
			}
			n.campaignService.DeleteCampaign(campaign.ID)

		}(&campaign)
	}

}

// checkContributionReminders checks if a campaign has contributions that are due for reminders
func (n *cronService) checkContributionReminders() {
	campaigns, err := n.campaignService.GetActiveCampaigns()
	if err != nil {
		return
	}
	for _, campaign := range campaigns {
		for _, contributor := range campaign.Contributors {
			if !contributor.HasPaid() {
				go n.notificationService.SendContributionReminder(&contributor, &campaign)
			}
		}
	}

}

// checkCampaignDeadline checks if a campaign is about to end and sends a reminder
func (n *cronService) checkCampaignDeadline() {
	campaigns, err := n.campaignService.GetNearEndCampaigns()
	if err != nil {
		return
	}
	for _, campaign := range campaigns {
		go n.notificationService.SendDeadlineReminder(&campaign)
	}
}

// Helper Functions ----------------------------------------------

func createJSONExport(data models.Campaign) (string, error) {
	fileName := fmt.Sprintf("campaign_export_%s_%s.json",
		data.ID,
		time.Now().UTC().Format("2006-01-02"))

	filePath := fmt.Sprintf("/tmp/%s", fileName)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("error writing JSON file: %w", err)
	}

	return filePath, nil
}
