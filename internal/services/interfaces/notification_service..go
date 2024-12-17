package interfaces

type NotificationService interface {
	// NotifyCampaignCreation(campaign *models.Campaign, contributors []models.Contributor) error
	// NotifyActivityAddition(activity *models.Activity, campaign *models.Campaign) error
	// NotifyContributorAdded(contributor *models.Contributor, campaign *models.Campaign) error
	// SendContributionReminder(contributor *models.Contributor, campaign *models.Campaign, dueDate time.Time) error
	// SendBulkContributionReminders(campaign *models.Campaign) error
	// NotifyActivityApprovalRequest(activity *models.Activity, requestedBy string) error
	// NotifyActivityApproval(activity *models.Activity, approvedBy string, approvalTime time.Time) error
	// NotifyActivityUpdate(activity *models.Activity, updateType string) error
	// NotifyCampaignUpdate(campaign *models.Campaign, updateType string) error
	// NotifyCampaignMilestone(campaign *models.Campaign, milestoneType string) error
	// SendDeadlineReminder(campaign *models.Campaign, deadline time.Time) error
	// NotifyPaymentReceived(contributor *models.Contributor, amount float64, campaign *models.Campaign) error
	// NotifyPaymentDue(contributor *models.Contributor, amount float64, dueDate time.Time) error
	// SendSystemNotification(notificationType string, message string, recipients []string) error
}
