package templates

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

func generateFile(fileName string) string {
	//TODO:implement better file path generation and handling
	return filepath.Join(config.BaseDir, "pkg", "email", "templates", fileName)
}

// Personal Email Templates
func AnalyticsReport(to []string, today, comparison interface{}, reportDate time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: fmt.Sprintf("Daily Analytics Report - %s", reportDate.Format("2006-01-02")),
		Path:    generateFile("personal/analytics_report.html"),
		Data: map[string]interface{}{
			"today":      today,
			"comparison": comparison,
			"date":       reportDate.Format("January 2, 2006"),
		},
	}
}

func Verification(to []string, name, verificationCode string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Email Verification - GoFund It",
		Path:    generateFile("personal/email_verification.html"),
		Data: map[string]interface{}{
			"verificationCode": verificationCode,
			"name":             name,
		},
	}
}

func ActivityApprovalRequest(to []string, campaignTitle, activityTitle, activitySubtitle string, activityCost float64, requestedBy string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Activity Approval Request - GoFund It",
		Path:    generateFile("personal/activity_approval_request.html"),
		Data: map[string]interface{}{
			"Campaign": map[string]string{
				"Title": campaignTitle,
			},
			"Activity": map[string]interface{}{
				"Title":    activityTitle,
				"Subtitle": activitySubtitle,
				"Cost":     activityCost,
			},
			"RequestedBy": requestedBy,
		},
	}
}

func CampaignCreated(to []string, campaignTitle, campaignDescription, campaignID, campaignKey string, contributors, activities []map[string]string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Created - GoFund It",
		Path:    generateFile("personal/campaign_created.html"),
		Data: map[string]interface{}{
			"title":        campaignTitle,
			"description":  campaignDescription,
			"id":           campaignID,
			"key":          campaignKey,
			"contributors": contributors,
			"activities":   activities,
		},
	}
}

func CampaignCleanUp(to []string, campaignName string, endDate time.Time, totalAmount float64, campaignDataPath string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:          to,
		Subject:     "Campaign Data Export - GoFund It",
		Path:        generateFile("general/campaign_cleanup.html"),
		Attachments: []string{campaignDataPath},
		Data: map[string]interface{}{
			"CampaignName": campaignName,
			"EndDate":      endDate.Format("January 2, 2006"),
			"TotalAmount":  totalAmount,
		},
	}
}

func CampaignDeadlineReminder(to []string, campaignTitle string, deadline time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Deadline Reminder - GoFund It",
		Path:    generateFile("personal/campaign_deadline.html"),
		Data: map[string]interface{}{
			"title":    campaignTitle,
			"deadline": deadline.Format("January 2, 2006"),
		},
	}
}

func ContributorAdded(to []string, name, campaignTitle, campaignID, campaignKey string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Contributor Added - GoFund It",
		Path:    generateFile("personal/contributor_added.html"),
		Data: map[string]interface{}{
			"name":  name,
			"title": campaignTitle,
			"id":    campaignID,
			"key":   campaignKey,
		},
	}
}

func ContributionReminder(to []string, name, campaignTitle string, dueDate time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Contribution Reminder - GoFund It",
		Path:    generateFile("personal/contribution_reminder.html"),
		Data: map[string]interface{}{
			"name":          name,
			"campaignTitle": campaignTitle,
			"dueDate":       dueDate.Format("January 2, 2006"),
		},
	}
}

func PaymentReceived(to []string, name string, amount float64, campaignTitle string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Payment Received - GoFund It",
		Path:    generateFile("personal/payment_received.html"),
		Data: map[string]interface{}{
			"name":          name,
			"amount":        amount,
			"campaignTitle": campaignTitle,
		},
	}
}

func PayoutRequired(to []string, campaignID string, payoutAmount float64, endDate, cleanupDate time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Payout Required - GoFund It",
		Path:    generateFile("personal/campaign_payout_required.html"),
		Data: map[string]interface{}{
			"payoutAmount": payoutAmount,
			"endDate":      endDate.Format("January 2, 2006"),
			"cleanupDate":  cleanupDate.Format("January 2, 2006"),
			"campaignId":   campaignID,
		},
	}
}

func CampaignEnded(to []string, campaignTitle string, contribution, totalRaised float64, endDate time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Ended - GoFund It",
		Path:    generateFile("personal/campaign_ended.html"),
		Data: map[string]interface{}{
			"campaignTitle": campaignTitle,
			"contribution":  contribution,
			"totalRaised":   totalRaised,
			"endDate":       endDate.Format("January 2, 2006"),
			"currentYear":   time.Now().Year(),
		},
	}
}

// General Email Templates

func ActivityAddedGeneral(to []string, campaignTitle, activityTitle, activitySubtitle string, activityCost float64) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "New Activity Added - GoFund It",
		Path:    generateFile("general/activity_added.html"),
		Data: map[string]interface{}{
			"campaignTitle":    campaignTitle,
			"activityTitle":    activityTitle,
			"activitySubtitle": activitySubtitle,
			"activityCost":     activityCost,
		},
	}
}

func ActivityApprovedGeneral(to []string, campaignTitle, activityTitle, activitySubtitle, approvedBy string, approvalTime time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Activity Approved - GoFund It",
		Path:    generateFile("general/activity_approved.html"),
		Data: map[string]interface{}{
			"campaignTitle":    campaignTitle,
			"activityTitle":    activityTitle,
			"activitySubtitle": activitySubtitle,
			"approvedBy":       approvedBy,
			"approvalTime":     approvalTime.Format("January 2, 2006, 15:04"),
		},
	}
}

func ActivityUpdateGeneral(to []string, campaignTitle, activityTitle, updateType string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Activity Update - GoFund It",
		Path:    generateFile("general/activity_updated.html"),
		Data: map[string]interface{}{
			"campaignTitle": campaignTitle,
			"activityTitle": activityTitle,
			"updateType":    updateType,
		},
	}
}

func CampaignMilestoneGeneral(to []string, campaignTitle, milestoneType string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Milestone - GoFund It",
		Path:    generateFile("general/campaign_milestone.html"),
		Data: map[string]interface{}{
			"title":         campaignTitle,
			"milestoneType": milestoneType,
		},
	}
}

func CampaignUpdatedGeneral(to []string, campaignTitle, updateType string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Updated - GoFund It",
		Path:    generateFile("general/campaign_updated.html"),
		Data: map[string]interface{}{
			"title":      campaignTitle,
			"updateType": updateType,
		},
	}
}

func CommentAddedGeneral(to []string, commenterName, commentContent, activityTitle, campaignTitle string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "New Comment on Activity - GoFund It",
		Path:    generateFile("general/comment_added.html"),
		Data: map[string]interface{}{
			"commenterName":  commenterName,
			"commentContent": commentContent,
			"activityTitle":  activityTitle,
			"campaignTitle":  campaignTitle,
		},
	}
}

func ContributorAddedGeneral(to []string, contributorName, contributorEmail string, contributorAmount float64, campaignTitle string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "New Contributor Added - GoFund It",
		Path:    generateFile("general/contributor_added.html"),
		Data: map[string]interface{}{
			"contributorName":   contributorName,
			"contributorEmail":  contributorEmail,
			"contributorAmount": contributorAmount,
			"campaignTitle":     campaignTitle,
		},
	}
}

func PaymentNotificationGeneral(to []string, contributorName string, amount float64, campaignTitle string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Payment Notification - GoFund It",
		Path:    generateFile("general/payment_notification.html"),
		Data: map[string]interface{}{
			"contributorName": contributorName,
			"amount":          amount,
			"campaignTitle":   campaignTitle,
		},
	}
}

func PayoutCollected(to []string, campaignTitle, ownerName string, payoutAmount float64, payoutDate time.Time) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "Campaign Payout Collected - GoFund It",
		Path:    generateFile("general/payout_collected.html"),
		Data: map[string]interface{}{
			"campaignTitle": campaignTitle,
			"ownerName":     ownerName,
			"payoutAmount":  payoutAmount,
			"payoutDate":    payoutDate.Format("January 2, 2006"),
			"currentYear":   time.Now().Year(),
		},
	}
}
func SystemNotificationGeneral(to []string, notificationType, message string) *email.EmailTemplate {
	return &email.EmailTemplate{
		To:      to,
		Subject: "System Notification - GoFund It",
		Path:    generateFile("general/system_notification.html"),
		Data: map[string]interface{}{
			"notificationType": notificationType,
			"message":          message,
		},
	}
}
