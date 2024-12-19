package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ActivityRepository interface {
	Create(activity *models.Activity) (models.Activity, error)
	Update(activity *models.Activity) error
	Delete(activity *models.Activity) error
	GetActivityByID(activityID uint) (models.Activity, error)
	GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error)
	GetActivityParticipants(activityID uint) ([]models.Contributor, error)
	UpdateActivity(activity *models.Activity) error
	AddContributorToActivity(activityID uint, contributorID uint) error
	RemoveContributorFromActivity(activityID uint, contributorID uint) error
}
