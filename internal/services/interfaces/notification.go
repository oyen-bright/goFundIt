package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type NotificationService interface {
	// Campaign notifications
	NotifyCampaignCleanUp(campaign *models.Campaign, data string) error
	NotifyCampaignCreation(campaign *models.Campaign) error
	NotifyCampaignUpdate(campaign *models.Campaign, updateType string) error
	NotifyCampaignMilestone(campaign *models.Campaign, milestoneType string) error

	// Activity notifications
	NotifyActivityAddition(activity *models.Activity, campaign *models.Campaign) error
	NotifyActivityApproved(activity *models.Activity, campaign *models.Campaign) error
	NotifyActivityApprovalRequest(activity *models.Activity, campaign *models.Campaign) error
	NotifyActivityUpdate(activity *models.Activity, campaign *models.Campaign) error

	//Comment notifications
	NotifyCommentAddition(comment *models.Comment, activityID *models.Activity) error

	// Contributor notifications
	NotifyContributorAdded(contributor *models.Contributor, campaign *models.Campaign) error
	NotifyPaymentReceived(contributor *models.Contributor, campaign *models.Campaign) error

	// Payout
	NotifyPayoutCollected(campaign *models.Campaign) error
	NotifyCampaignPayoutRequired(campaign *models.Campaign) error

	// Reminder notifications
	SendContributionReminder(contributor *models.Contributor, campaign *models.Campaign) error
	SendDeadlineReminder(campaign *models.Campaign) error

	// System notifications
	SendSystemNotification(notificationType string, message string) error
}
