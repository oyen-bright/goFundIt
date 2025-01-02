package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ActivityService interface {
	CreateActivity(activity models.Activity, userHandle, campaignId, key string) (models.Activity, error)
	UpdateActivity(activity *models.Activity, userHandle string) error
	DeleteActivityByID(activityID uint, campaignID, userHandle string) error

	GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error)
	GetActivityByID(activityID uint, campaignID string) (models.Activity, error)
	GetParticipants(activityID uint, campaignId, key string) ([]models.Contributor, error)

	OptInContributor(campaignID, userEmail, key string, activityID, contributorID uint) error
	OptOutContributor(campaignID, userEmail, key string, activityID, contributorID uint) error

	ApproveActivity(activityID uint, userHandle, key string) (*models.Activity, error)
}
