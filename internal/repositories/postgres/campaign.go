package postgress

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type campaignRepository struct {
	db *gorm.DB
}

// NewCampaignRepository creates a new campaign repository instance
func NewCampaignRepository(db *gorm.DB) interfaces.CampaignRepository {
	return &campaignRepository{db: db}
}

// Create creates a new campaign
func (r *campaignRepository) Create(campaign *models.Campaign) (models.Campaign, error) {
	if err := r.db.Create(campaign).Error; err != nil {
		return models.Campaign{}, err
	}
	return *campaign, nil
}

// Update updates a campaign
func (r *campaignRepository) Update(campaign *models.Campaign) (models.Campaign, error) {
	if err := r.db.Save(campaign).Error; err != nil {
		return models.Campaign{}, err
	}
	return *campaign, nil
}

// Delete deletes a campaign
func (r *campaignRepository) Delete(campaignID string) error {
	return r.db.Where("id = ?", campaignID).Delete(&models.Campaign{}).Error
}

// TODO: Redundant ? GetByIDWithSelectedData
func (r *campaignRepository) GetByID(id string) (models.Campaign, error) {
	var campaign models.Campaign

	query := r.db.Where("id = ?", id)
	query = query.Preload("Images").Preload("Activities.Contributors").Preload("Activities").Preload("Contributors.Payment").Preload("Contributors").Preload("Payout").Preload("CreatedBy")
	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

// TODO:payout createdby not loading
// GetByIDWithSelectedData fetches a campaign by ID with selected data preloaded
func (r *campaignRepository) GetByIDWithSelectedData(id string, options models.PreloadOption) (models.Campaign, error) {
	var campaign models.Campaign
	query := r.db.Where("id = ?", id)

	// Base campaign query
	if options.Images {
		query = query.Preload("Images")
	}

	if options.Payout {
		query = query.Preload("Payout")
	}

	if options.Activities {
		query = query.Preload("Activities", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		})

		if options.ActivitiesContributors {
			query = query.Preload("Activities.Contributors")
		}

		if options.ActivitiesComments {
			// Load comments with nested replies and their creators
			query = query.
				Preload("Activities.Comments", func(db *gorm.DB) *gorm.DB {
					return db.Order("created_at DESC").Where("parent_id IS NULL")
				}).
				Preload("Activities.Comments.CreatedBy").
				Preload("Activities.Comments.Replies").
				Preload("Activities.Comments.Replies.CreatedBy").
				Preload("Activities.Comments.Replies.Replies").
				Preload("Activities.Comments.Replies.Replies.CreatedBy")
		}

		// Always preload Activity CreatedBy if Activities are loaded
		query = query.Preload("Activities.CreatedBy")
	}

	if options.Contributors {
		query = query.Preload("Contributors", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).Preload("Contributors.Payment")

		if options.ContributorsActivities {
			query = query.Preload("Contributors.Activities")
		}
	}

	query = query.Preload("CreatedBy").Preload("Payout")
	err := query.First(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// GetByHandle fetches a campaign by creator's handle
func (r *campaignRepository) GetByHandle(handle string) (models.Campaign, error) {
	var campaign models.Campaign
	query := r.db.Where("created_by_handle = ?", handle)
	query = query.Preload("Images").Preload("Activities").Preload("Contributors")
	query = query.Preload("CreatedBy")
	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

// GetExpiredCampaigns fetches all expired campaigns
func (r *campaignRepository) GetExpiredCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	query := r.db.Where("end_date <= ?", time.Now().UTC())
	query = query.Preload("Contributors.Payment").Preload("Contributors").Preload("CreatedBy")
	err := query.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

// GetActiveCampaigns fetches all active campaigns
func (r *campaignRepository) GetActiveCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	query := r.db.Where("end_date > ?", time.Now().UTC())
	query = query.Preload("Contributors.Payment").Preload("Contributors")
	query = query.Preload("CreatedBy")
	err := query.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

// GetNearEndCampaigns fetches all campaigns that are near end in 3 days
func (r *campaignRepository) GetNearEndCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	now := time.Now().UTC()
	threeDaysFromNow := now.AddDate(0, 0, 3)

	query := r.db.Where("end_date BETWEEN ? AND ?", now, threeDaysFromNow)
	query = query.Or("end_date = ?", now.AddDate(0, 0, 1))
	query = query.Preload("Contributors.Payment").Preload("Contributors")
	query = query.Preload("CreatedBy")

	err := query.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}
