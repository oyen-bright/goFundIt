package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type ActivityRepository interface {
	Create(activity *models.Activity) (models.Activity, error)
	Update(activity *models.Activity) error
	Delete(activity *models.Activity) error
	GetByID(activityID uint) (models.Activity, error)
	GetByCampaignID(campaignID string) ([]models.Activity, error)
	GetParticipants(activityID uint) ([]models.Contributor, error)
	Save(activity *models.Activity) error
	AddContributor(activityID uint, contributorID uint) error
	RemoveContributor(activityID uint, contributorID uint) error
}
