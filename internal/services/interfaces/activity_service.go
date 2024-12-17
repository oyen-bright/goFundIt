package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ActivityService interface {
	CreateActivity(activity models.Activity, userHandle, campaignId string) (models.Activity, error)
	GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error)
	GetActivityByID(activityID uint, campaignID string) (models.Activity, error)
	UpdateActivity(activity *models.Activity, userHandle string) error
	DeleteActivityByID(activityID uint, campaignID, userHandle string) error
	// JoinActivity(contributorID uint, activityID uint) error
	// LeaveActivity(contributorID uint, activityID uint) error
	// GetActivityParticipants(activityID uint) ([]contributor.Contributor, error)
}
