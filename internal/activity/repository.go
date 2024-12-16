package activity

import (
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *Activity) error
	Update(activity *Activity) error
	Delete(activity *Activity) error
	GetActivitiesByCampaignID(campaignID string) ([]Activity, error)
	GetActivityByID(activityID uint) (Activity, error)
	GetActivityParticipants(activityID uint) ([]contributor.Contributor, error)
}

type activityRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) Update(activity *Activity) error {
	return r.db.Model(&Activity{}).Where("id = ?", activity.ID).Updates(
		map[string]interface{}{
			"title":       activity.Title,
			"subtitle":    activity.Subtitle,
			"imageUrl":    activity.ImageUrl,
			"isMandatory": activity.IsMandatory,
			"cost":        activity.Cost,
			"isApproved":  activity.IsApproved,
		}).Error
}

func (r *activityRepository) Delete(activity *Activity) error {
	return r.db.Delete(activity).Error
}
func (r *activityRepository) GetActivitiesByCampaignID(campaignID string) ([]Activity, error) {
	var activities []Activity
	err := r.db.Preload("Contributors").Where("campaign_id = ?", campaignID).Find(&activities).Error
	return activities, err
}

func (r *activityRepository) GetActivityByID(activityID uint) (Activity, error) {
	var activity Activity
	err := r.db.Preload("Contributors").First(&activity, activityID).Error
	return activity, err
}

func (r *activityRepository) UpdateActivity(activity *Activity) error {
	return r.db.Save(activity).Error
}

func (r *activityRepository) DeleteActivity(activityID uint) error {
	return r.db.Delete(&Activity{}, activityID).Error
}

func (r *activityRepository) GetActivityParticipants(activityID uint) ([]contributor.Contributor, error) {
	var participants []contributor.Contributor
	//TODO:"fix query"
	err := r.db.Joins("JOIN activity_contributors ON activity_contributors.contributor_id = contributors.id").
		Where("activity_contributors.activity_id = ?", activityID).Find(&participants).Error
	return participants, err
}
