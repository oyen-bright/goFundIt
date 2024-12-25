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

func NewCampaignRepository(db *gorm.DB) interfaces.CampaignRepository {
	return &campaignRepository{db: db}
}

// TODO: redundant GetByIDWithSelectedData
func (r *campaignRepository) GetByID(id string, preload bool) (models.Campaign, error) {
	var campaign models.Campaign

	query := r.db.Where("id = ?", id)

	if preload {
		query = query.Preload("Images").Preload("Activities.Contributors").Preload("Activities").Preload("Contributors.Payment").Preload("Contributors").Preload("CreatedBy")
	}

	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

// TODO: redundant GetByIDWithSelectedData
func (r *campaignRepository) GetByIDWithContributors(id string) (models.Campaign, error) {
	var campaign models.Campaign

	query := r.db.Where("id = ?", id)

	query = query.Preload("Contributors.Activities").Preload("Contributors.Payment").Preload("Contributors").Preload("CreatedBy")

	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

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

	query = query.Preload("CreatedBy")

	err := query.First(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepository) GetByCreatorHandle(handle string, preload bool) (models.Campaign, error) {
	var campaign models.Campaign

	query := r.db.Where("created_by_handle = ?", handle)

	if preload {
		query = query.Preload("Images").Preload("Activities").Preload("Contributors")
	}
	query = query.Preload("CreatedBy")
	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

func (r *campaignRepository) Create(campaign *models.Campaign) (models.Campaign, error) {
	if err := r.db.Create(campaign).Error; err != nil {
		return models.Campaign{}, err
	}
	return *campaign, nil
}
func (r *campaignRepository) Update(campaign *models.Campaign) (models.Campaign, error) {
	if err := r.db.Save(campaign).Error; err != nil {
		return models.Campaign{}, err
	}
	return *campaign, nil
}

func (r *campaignRepository) Delete(campaignID string) error {
	return r.db.Where("id = ?", campaignID).Delete(&models.Campaign{}).Error
}

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

// func (r *campaignRepository) Update(campaign *models.Campaign) (models.Campaign, error) {
// 	err := r.db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Save(campaign).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Model(campaign).Association("Activities").Replace(campaign.Activities); err != nil {
// 			return err
// 		}

// 		if err := tx.Model(campaign).Association("Images").Replace(campaign.Images); err != nil {
// 			return err
// 		}

// 		if err := tx.Model(campaign).Association("Contributors").Replace(campaign.Contributors); err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return models.Campaign{}, err
// 	}
// 	return *campaign, nil
// }

// func (r *campaignRepository) Delete(campaignID string) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.Activity{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.CampaignImage{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.Contributor{}).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Where("id = ?", campaignID).Delete(&models.Campaign{}).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

//TODO: implement preload builde

// // PreloadOptionBuilder provides methods to build PreloadOption
// type PreloadOptionBuilder struct {
//     options PreloadOption
// }

// // NewPreloadOptionBuilder creates a new PreloadOptionBuilder
// func NewPreloadOptionBuilder() *PreloadOptionBuilder {
//     return &PreloadOptionBuilder{}
// }

// // WithImages includes Images in the preload
// func (b *PreloadOptionBuilder) WithImages() *PreloadOptionBuilder {
//     b.options.Images = true
//     return b
// }

// // WithPayout includes Payout in the preload
// func (b *PreloadOptionBuilder) WithPayout() *PreloadOptionBuilder {
//     b.options.Payout = true
//     return b
// }

// // WithActivities includes Activities and their CreatedBy in the preload
// func (b *PreloadOptionBuilder) WithActivities() *PreloadOptionBuilder {
//     b.options.Activities = true
//     return b
// }

// // WithActivityContributors includes Activity Contributors in the preload
// func (b *PreloadOptionBuilder) WithActivityContributors() *PreloadOptionBuilder {
//     b.options.Activities = true
//     b.options.ActivitiesContributore = true
//     return b
// }

// // WithActivityComments includes Activity Comments in the preload
// func (b *PreloadOptionBuilder) WithActivityComments() *PreloadOptionBuilder {
//     b.options.Activities = true
//     b.options.ActiviitesComments = true
//     return b
// }

// // WithContributors includes Contributors in the preload
// func (b *PreloadOptionBuilder) WithContributors() *PreloadOptionBuilder {
//     b.options.Contributors = true
//     return b
// }

// // WithContributorActivities includes Contributor Activities in the preload
// func (b *PreloadOptionBuilder) WithContributorActivities() *PreloadOptionBuilder {
//     b.options.Contributors = true
//     b.options.ContributorsActivities = true
//     return b
// }

// // WithCreatedBy includes CreatedBy in the preload
// func (b *PreloadOptionBuilder) WithCreatedBy() *PreloadOptionBuilder {
//     b.options.CreatedBy = true
//     return b
// }

// // Build creates the PreloadOption
// func (b *PreloadOptionBuilder) Build() PreloadOption {
//     return b.options
// }
