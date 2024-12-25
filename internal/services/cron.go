package services

import (
	"fmt"
	"os"
	"time"

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

// StartCronJobs implements interfaces.NotificationService.
func (n *cronService) StartCronJobs() error {

	// Daily cleanup at midnight UTC
	_, err := n.cron.AddFunc("0 0 0 * * *", func() {
		n.cleanUpExpiredCampaign()
	})
	if err != nil {
		return fmt.Errorf("failed to schedule cleanup job: %w", err)
	}

	// Check contribution reminders every 3 days
	_, err = n.cron.AddFunc("0 0 */3 * *", func() {
		n.checkContributionReminders()
	})
	if err != nil {
		return fmt.Errorf("failed to schedule contribution reminders job: %w", err)
	}

	// Check campaign deadlines every day
	_, err = n.cron.AddFunc("0 0 * * *", func() {
		n.checkCampaignDeadline()
	})
	if err != nil {
		return fmt.Errorf("failed to schedule campaign deadline reminders job: %w", err)
	}

	return nil

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
