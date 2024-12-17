package postgress

import (
	"fmt"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepo(db *gorm.DB) interfaces.ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *models.Activity) (models.Activity, error) {
	err := r.db.Create(activity).Error
	return *activity, err
}

func (r *activityRepository) Update(activity *models.Activity) error {
	return r.db.Model(&models.Activity{}).Where("id = ?", activity.ID).Updates(
		map[string]interface{}{
			"title":        activity.Title,
			"subtitle":     activity.Subtitle,
			"image_url":    activity.ImageUrl,
			"is_mandatory": activity.IsMandatory,
			"cost":         activity.Cost,
			"is_approved":  activity.IsApproved,
		}).Error
}

func (r *activityRepository) Delete(activity *models.Activity) error {
	return r.db.Delete(activity).Error
}

func (r *activityRepository) GetActivityByID(activityID uint) (models.Activity, error) {
	var activity models.Activity
	err := r.db.Preload("Contributors").First(&activity, activityID).Error

	fmt.Println(activity)
	return activity, err
}

func (r *activityRepository) GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Preload("Contributors").Where("campaign_id = ?", campaignID).Find(&activities).Error
	return activities, err
}

func (r *activityRepository) UpdateActivity(activity *models.Activity) error {
	return r.db.Save(activity).Error
}

func (r *activityRepository) DeleteActivity(activityID uint) error {
	return r.db.Delete(&models.Activity{}, activityID).Error
}

func (r *activityRepository) GetActivityParticipants(activityID uint) ([]models.Contributor, error) {
	var participants []models.Contributor
	//TODO:"fix query"
	err := r.db.Joins("JOIN activity_contributors ON activity_contributors.contributor_id = contributors.id").
		Where("activity_contributors.activity_id = ?", activityID).Find(&participants).Error
	return participants, err
}
