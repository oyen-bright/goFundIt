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

// Save saves changes to an existing activity
func (r *activityRepository) Save(activity *models.Activity) error {
	return r.db.Save(activity).Error
}

// Delete removes an activity from the database
func (r *activityRepository) Delete(activity *models.Activity) error {
	return r.db.Delete(activity).Error
}

// GetActivityByID retrieves a single activity by its ID with contributors
func (r *activityRepository) GetByID(activityID uint) (models.Activity, error) {
	var activity models.Activity
	err := r.db.Preload("Contributors").First(&activity, activityID).Error

	fmt.Println(activity)
	return activity, err
}

// RemoveContributorFromActivity removes a contributor from an activity
func (r *activityRepository) RemoveContributor(activityID uint, contributorID uint) error {
	return r.db.Table("activities_contributors").
		Where("activity_id = ? AND contributor_id = ?", activityID, contributorID).
		Delete(nil).Error
}

// AddContributorToActivity adds a contributor to an activity
func (r *activityRepository) AddContributor(activityID uint, contributorID uint) error {
	return r.db.Table("activities_contributors").Create(map[string]interface{}{
		"activity_id":    activityID,
		"contributor_id": contributorID,
	}).Error
}

// GetActivitiesByCampaignID fetches all activities for a specific campaign
func (r *activityRepository) GetByCampaignID(campaignID string) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Preload("Contributors").Where("campaign_id = ?", campaignID).Find(&activities).Error
	return activities, err
}

// GetActivityParticipants retrieves all contributors for a specific activity
func (r *activityRepository) GetParticipants(activityID uint) ([]models.Contributor, error) {
	var participants []models.Contributor
	err := r.db.Table("contributors").
		Joins("JOIN activities_contributors ON activities_contributors.contributor_id = contributors.id").
		Where("activities_contributors.activity_id = ?", activityID).Find(&participants).Error
	return participants, err
}
