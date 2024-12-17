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

// Create adds a new activity to the database
func (r *activityRepository) Create(activity *models.Activity) (models.Activity, error) {
	err := r.db.Create(activity).Error
	return *activity, err
}

// Update modifies an existing activity's details
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

// Delete removes an activity from the database
func (r *activityRepository) Delete(activity *models.Activity) error {
	return r.db.Delete(activity).Error
}

// GetActivityByID retrieves a single activity by its ID with contributors
func (r *activityRepository) GetActivityByID(activityID uint) (models.Activity, error) {
	var activity models.Activity
	err := r.db.Preload("Contributors").First(&activity, activityID).Error

	fmt.Println(activity)
	return activity, err
}

// GetActivitiesByCampaignID fetches all activities for a specific campaign
func (r *activityRepository) GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Preload("Contributors").Where("campaign_id = ?", campaignID).Find(&activities).Error
	return activities, err
}

// UpdateActivity saves changes to an existing activity
func (r *activityRepository) UpdateActivity(activity *models.Activity) error {
	return r.db.Save(activity).Error
}

// DeleteActivity removes an activity by its ID
func (r *activityRepository) DeleteActivity(activityID uint) error {
	return r.db.Delete(&models.Activity{}, activityID).Error
}

// GetActivityParticipants retrieves all contributors for a specific activity
func (r *activityRepository) GetActivityParticipants(activityID uint) ([]models.Contributor, error) {
	var participants []models.Contributor
	//TODO:"fix query"
	err := r.db.Joins("JOIN activity_contributors ON activity_contributors.contributor_id = contributors.id").
		Where("activity_contributors.activity_id = ?", activityID).Find(&participants).Error
	return participants, err
}
