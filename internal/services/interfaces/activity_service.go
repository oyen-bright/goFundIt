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
	OptInContributor(campaignID, userEmail string, activityID, contributorID uint) error
	OptOutContributor(campaignID, userEmail string, activityID, contributorID uint) error
	GetParticipants(activityID uint, campaignId string) ([]models.Contributor, error)
	ApproveActivity(activityID uint, userHandle string) (*models.Activity, error)
}
