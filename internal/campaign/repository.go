package campaign

import (
	"github.com/oyen-bright/goFundIt/internal/activity"
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"gorm.io/gorm"
)

type CampaignRepository interface {
	Create(campaign *Campaign) (Campaign, error)
	Update(campaign *Campaign) (Campaign, error)
	Delete(campaignID string) error

	GetByID(id string, preload bool) (Campaign, error)
	GetByCreatorHandle(handle string, preload bool) (Campaign, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) GetByID(id string, preload bool) (Campaign, error) {
	var campaign Campaign
	query := r.db.Where("id = ?", id)

	if preload {
		query = query.Preload("Images").Preload("Activities").Preload("Contributors")
	}

	err := query.First(&campaign).Error
	if err != nil {
		return Campaign{}, err
	}
	return campaign, nil
}

func (r *campaignRepository) GetByCreatorHandle(handle string, preload bool) (Campaign, error) {
	var campaign Campaign
	query := r.db.Where("created_by_handle = ?", handle)

	if preload {
		query = query.Preload("Images").Preload("Activities").Preload("Contributors")
	}

	err := query.First(&campaign).Error
	if err != nil {
		return Campaign{}, err
	}
	return campaign, nil
}

func (r *campaignRepository) Create(campaign *Campaign) (Campaign, error) {
	if err := r.db.Create(campaign).Error; err != nil {
		return Campaign{}, err
	}
	return *campaign, nil
}

func (r *campaignRepository) Update(campaign *Campaign) (Campaign, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(campaign).Error; err != nil {
			return err
		}

		if err := tx.Model(campaign).Association("Activities").Replace(campaign.Activities); err != nil {
			return err
		}

		if err := tx.Model(campaign).Association("Images").Replace(campaign.Images); err != nil {
			return err
		}

		if err := tx.Model(campaign).Association("Contributors").Replace(campaign.Contributors); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Campaign{}, err
	}
	return *campaign, nil
}

func (r *campaignRepository) Delete(campaignID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("campaign_id = ?", campaignID).Delete(&activity.Activity{}).Error; err != nil {
			return err
		}

		if err := tx.Where("campaign_id = ?", campaignID).Delete(&CampaignImage{}).Error; err != nil {
			return err
		}

		if err := tx.Where("campaign_id = ?", campaignID).Delete(&contributor.Contributor{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", campaignID).Delete(&Campaign{}).Error; err != nil {
			return err
		}

		return nil
	})
}
