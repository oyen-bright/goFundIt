package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ActivityService interface {
	CreateActivity(activity models.Activity, userHandle, campaignId string) (models.Activity, error)
	// GetActivitiesByCampaignID(campaignID string) ([]Activity, error)
	// GetActivityByID(activityID uint) (Activity, error)
	// UpdateActivity(activity *Activity) error
	// DeleteActivity(activityID uint) error
	// JoinActivity(contributorID uint, activityID uint) error
	// LeaveActivity(contributorID uint, activityID uint) error
	// GetActivityParticipants(activityID uint) ([]contributor.Contributor, error)
}
