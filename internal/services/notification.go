package services

import (
	"context"
	"fmt"
	"log"

	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	"github.com/oyen-bright/goFundIt/pkg/fcm"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

//TODO:FCM for contributors

type emailNotifier struct {
	client email.Emailer
	logger logger.Logger
}

type fcmNotifier struct {
	client *fcm.Client
	logger logger.Logger
}

type notificationService struct {
	emailer     emailNotifier
	fcmNotifier fcmNotifier
	authService services.AuthService
	logger      logger.Logger
}

// notificationService implements interfaces.NotificationService.
func NewNotificationService(emailer email.Emailer, authService services.AuthService, fcmClient *fcm.Client, logger logger.Logger) services.NotificationService {
	return &notificationService{
		emailer:     emailNotifier{client: emailer, logger: logger},
		fcmNotifier: fcmNotifier{client: fcmClient, logger: logger},
		authService: authService,
		logger:      logger,
	}
}

func (e *emailNotifier) send(template *email.EmailTemplate) error {
	if err := e.client.SendEmailTemplate(*template); err != nil {
		e.logger.Error(err, "Error sending email: "+template.Name+err.Error(), nil)
		return err
	}
	return nil
}

func (f *fcmNotifier) send(data fcm.NotificationData, tokens []string) error {
	ctx := context.Background()
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) == 1 {

		err := f.client.SendNotification(ctx, tokens[0], data)
		if err != nil {
			f.logger.Error(err, "Error sending FCM message: "+data.Title+err.Error(), nil)
		}
		return err
	}

	err := f.client.SendMulticastNotification(ctx, tokens, data)
	if err != nil {
		f.logger.Error(err, "Error sending FCM multicast message: "+data.Title+err.Error(), nil)
	}
	return err

}

// ====== Activity Notifications ======

// NotifyActivityAddition sends an email to all activities of a campaign when a new activity is added.
func (n *notificationService) NotifyActivityAddition(activity *models.Activity, campaign *models.Campaign) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	activityAdded := emailTemplates.ActivityAddedGeneral(contributorsEmails, campaign.ID, activity.Title, activity.Subtitle, activity.Cost)

	return n.emailer.send(activityAdded)
}

// NotifyActivityApproval implements interfaces.NotificationService.
func (n *notificationService) NotifyActivityApproved(activity *models.Activity, campaign *models.Campaign) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	activityApprovedTemplate := emailTemplates.ActivityApprovedGeneral(contributorsEmails, campaign.ID, activity.Title, activity.Subtitle, campaign.CreatedBy.Email, activity.UpdatedAt)
	return n.emailer.send(activityApprovedTemplate)
}

// NotifyActivityApprovalRequest implements interfaces.NotificationService.
func (n *notificationService) NotifyActivityApprovalRequest(activity *models.Activity, campaign *models.Campaign) error {
	activityApprovalRequest := emailTemplates.ActivityApprovalRequest([]string{campaign.CreatedBy.Email}, campaign.ID, activity.Title, activity.Subtitle, activity.Cost, activity.CreatedBy.Email)
	userFCMToken := campaign.CreatedBy.FCMToken
	if userFCMToken != nil {
		n.fcmNotifier.send(fcm.NotificationData{
			Title: "New Activity Approval Request",
			Body:  fmt.Sprintf("A new activity  approval request has been added to %s", campaign.ID),
		}, []string{*userFCMToken})
	}
	return n.emailer.send(activityApprovalRequest)
}

// NotifyActivityUpdate implements interfaces.NotificationService.
func (n *notificationService) NotifyActivityUpdate(activity *models.Activity, campaign *models.Campaign) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	activityUpdate := emailTemplates.ActivityUpdateGeneral(contributorsEmails, campaign.ID, activity.Title, "details updated")
	return n.emailer.send(activityUpdate)
}

// ====== Campaign Notifications ======

// NotifyCampaignCreation implements interfaces.NotificationService.
func (n *notificationService) NotifyCampaignCreation(campaign *models.Campaign) error {
	contributorsNameEmail := getContributorNameEmail(campaign.Contributors)
	activitiesTitleSubtitle := getActivityTitleSubtitle(campaign.Activities)

	log.Println(contributorsNameEmail)
	log.Println(campaign)
	//Send email to campaign creator
	campaignCreatedCampaignCreator := emailTemplates.CampaignCreated([]string{campaign.CreatedBy.Email}, campaign.Title, campaign.Description, campaign.ID, campaign.Key, contributorsNameEmail, activitiesTitleSubtitle)
	err := n.emailer.send(campaignCreatedCampaignCreator)

	//send email to campaign contributors
	for _, contributor := range contributorsNameEmail {
		go func(contributor map[string]string) {
			contributorAddedTemplate := emailTemplates.ContributorAdded([]string{contributor["email"]}, contributor["name"], campaign.Title, campaign.ID, campaign.Key)
			err = n.emailer.send(contributorAddedTemplate)
		}(contributor)
	}
	return err
}

// NotifyCampaignMilestone implements interfaces.NotificationService.
func (n *notificationService) NotifyCampaignMilestone(campaign *models.Campaign, milestoneType string) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	contributorsEmails = append(contributorsEmails, campaign.CreatedBy.Email)

	userFCMToken := campaign.CreatedBy.FCMToken
	if userFCMToken != nil {
		n.fcmNotifier.send(fcm.NotificationData{
			Title: "New Campaign Milestone Reached",
			Body:  fmt.Sprintf("A new campaign milestone   %s", milestoneType),
		}, []string{*userFCMToken})
	}

	milestoneTemplate := emailTemplates.CampaignMilestoneGeneral(contributorsEmails, campaign.ID, milestoneType)
	return n.emailer.send(milestoneTemplate)
}

// NotifyCampaignUpdate implements interfaces.NotificationService.
func (n *notificationService) NotifyCampaignUpdate(campaign *models.Campaign, updateType string) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	contributorsEmails = append(contributorsEmails, campaign.CreatedBy.Email)

	campaignUpdateTemplate := emailTemplates.CampaignUpdatedGeneral(contributorsEmails, campaign.ID, updateType)
	return n.emailer.send(campaignUpdateTemplate)
}

// ====== Contributor Notifications ======

// NotifyContributorAdded implements interfaces.NotificationService.
func (n *notificationService) NotifyContributorAdded(contributor *models.Contributor, campaign *models.Campaign) error {
	contributorEmails := getContributorEmails(campaign.Contributors)
	contributorAddedTemplate := emailTemplates.ContributorAdded([]string{contributor.Email}, contributor.Name, campaign.Title, campaign.ID, campaign.Key)
	err := n.emailer.send(contributorAddedTemplate)
	if err != nil {
		return err
	}
	contributorAddedTemplate = emailTemplates.ContributorAddedGeneral(contributorEmails, contributor.Name, contributor.Email, contributor.Amount, campaign.ID)
	err = n.emailer.send(contributorAddedTemplate)
	return err
}

// ====== Payment and Payout Notifications ======

// NotifyPaymentReceived implements interfaces.NotificationService.
func (n *notificationService) NotifyPaymentReceived(contributor *models.Contributor, campaign *models.Campaign) error {
	paymentReceivedTemplate := emailTemplates.PaymentReceived([]string{contributor.Email}, contributor.Name, contributor.Payment.Amount, campaign.ID)

	userFCMToken := campaign.CreatedBy.FCMToken
	if userFCMToken != nil {
		n.fcmNotifier.send(fcm.NotificationData{
			Title: "New Payment Reached",
			Body:  fmt.Sprintf("A new payment received for campaign %s, from %s", campaign.ID, contributor.Email),
		}, []string{*userFCMToken})
	}
	return n.emailer.send(paymentReceivedTemplate)
}

// NotifyPayoutCollected implements interfaces.NotificationService.
func (n *notificationService) NotifyPayoutCollected(campaign *models.Campaign) error {
	contributorsEmails := getContributorEmails(campaign.Contributors)
	contributorsEmails = append(contributorsEmails, campaign.CreatedBy.Email)

	payoutCollectedTemplate := emailTemplates.PayoutCollected(contributorsEmails, campaign.ID, campaign.CreatedBy.Email, campaign.GetPayoutAmount(), campaign.Payout.UpdatedAt)
	return n.emailer.send(payoutCollectedTemplate)
}

// NotifyCampaignPayoutRequired implements interfaces.NotificationService.
func (n *notificationService) NotifyCampaignPayoutRequired(campaign *models.Campaign) error {
	payoutRequired := emailTemplates.PayoutRequired([]string{campaign.CreatedBy.Email}, campaign.ID, campaign.GetPayoutAmount(), campaign.EndDate, campaign.EndDate)

	userFCMToken := campaign.CreatedBy.FCMToken
	if userFCMToken != nil {
		n.fcmNotifier.send(fcm.NotificationData{
			Title: "Action Required: Campaign Payout",
			Body:  fmt.Sprintf("The campaign '%s' has reached its target. Complete the payout process by %s.", campaign.Title, campaign.EndDate.Format("02 Jan 2006")),
		}, []string{*userFCMToken})
	}
	return n.emailer.send(payoutRequired)
}

// ====== Reminder Notifications ======

// SendDeadlineReminder implements interfaces.NotificationService.
func (n *notificationService) SendDeadlineReminder(campaign *models.Campaign) error {
	payoutRequired := emailTemplates.CampaignDeadlineReminder([]string{campaign.CreatedBy.Email}, campaign.Title, campaign.EndDate)
	return n.emailer.send(payoutRequired)
}

// SendContributionReminder implements interfaces.NotificationService.
func (n *notificationService) SendContributionReminder(contributor *models.Contributor, campaign *models.Campaign) error {
	contributionReminder := emailTemplates.ContributionReminder([]string{contributor.Email}, contributor.Name, campaign.Title, campaign.EndDate)
	return n.emailer.send(contributionReminder)
}

// ====== System and Cleanup Notifications ======

// SendSystemNotification implements interfaces.NotificationService.
func (n *notificationService) SendSystemNotification(notificationType string, message string) error {
	allUsers, err := n.authService.GetAllUser()
	if err != nil {
		return err
	}
	allUsersEmail := make([]string, len(allUsers))
	for i, user := range allUsers {
		allUsersEmail[i] = user.Email
	}

	systemNotificationTemplate := emailTemplates.SystemNotificationGeneral(allUsersEmail, notificationType, message)
	return n.emailer.send(systemNotificationTemplate)
}

// NotifyCampaignCleanUp implements interfaces.NotificationService.
func (n *notificationService) NotifyCampaignCleanUp(campaign *models.Campaign, data string) error {
	campaignCleanUp := emailTemplates.CampaignCleanUp([]string{campaign.CreatedBy.Email}, campaign.ID, campaign.EndDate, campaign.GetPayoutAmount(), data)

	userFCMToken := campaign.CreatedBy.FCMToken
	if userFCMToken != nil {
		n.fcmNotifier.send(fcm.NotificationData{
			Title: "Campaign Data Cleanup Notification",
			Body:  fmt.Sprintf("All data related to the campaign '%s' has been cleaned and sent to your email. It is now removed from our system.", campaign.Title),
		}, []string{*userFCMToken})
	}
	return n.emailer.send(campaignCleanUp)
}

// ====== Comment Notifications ======

// NotifyCommentAddition implements interfaces.NotificationService.
func (n *notificationService) NotifyCommentAddition(comment *models.Comment, activity *models.Activity) error {
	contributorsEmails := getContributorEmails(activity.Contributors)
	contributorsEmails = append(contributorsEmails, activity.CreatedBy.Email)

	if len(contributorsEmails) == 0 {
		return nil
	}

	commentAddedTemplate := emailTemplates.CommentAddedGeneral(contributorsEmails, comment.CreatedBy.Handle, comment.Content, activity.Title, activity.CampaignID)
	return n.emailer.send(commentAddedTemplate)
}

// Helper Functions --------------------------------------------------

// getContributorEmails returns a list of emails from a list of contributors
func getContributorEmails(contributors []models.Contributor) []string {
	emails := make([]string, len(contributors))
	for i, contributor := range contributors {
		emails[i] = contributor.Email
	}
	return emails
}

// GetContributorNameEmail returns a
func getContributorNameEmail(contributors []models.Contributor) []map[string]string {
	data := make([]map[string]string, len(contributors))

	for i, contributor := range contributors {
		data[i] = map[string]string{
			"name":  contributor.Name,
			"email": contributor.Email,
		}
	}
	return data
}

// GetActivityTitleSubtitle
func getActivityTitleSubtitle(activities []models.Activity) []map[string]string {
	data := make([]map[string]string, len(activities))

	for i, activity := range activities {
		data[i] = map[string]string{
			"title":    activity.Title,
			"subtitle": activity.Subtitle,
		}
	}
	return data
}
