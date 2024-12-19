package postgress

import (
	"log"

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

func (r *campaignRepository) GetByID(id string, preload bool) (models.Campaign, error) {
	var campaign models.Campaign

	log.Println(id)
	query := r.db.Where("id = ?", id)

	if preload {
		query = query.Preload("Images").Preload("Activities.Contributors").Preload("Activities").Preload("Contributors").Preload("CreatedBy")
	}

	err := query.First(&campaign).Error
	if err != nil {
		return models.Campaign{}, err
	}
	log.Println(len(campaign.Activities))
	return campaign, nil
}

func (r *campaignRepository) GetByCreatorHandle(handle string, preload bool) (models.Campaign, error) {
	var campaign models.Campaign

	query := r.db.Where("created_by_handle = ?", handle)

	if preload {
		query = query.Preload("Images").Preload("Activities").Preload("Contributors")
	}

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
		return models.Campaign{}, err
	}
	return *campaign, nil
}

func (r *campaignRepository) Delete(campaignID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.Activity{}).Error; err != nil {
			return err
		}

		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.CampaignImage{}).Error; err != nil {
			return err
		}

		if err := tx.Where("campaign_id = ?", campaignID).Delete(&models.Contributor{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", campaignID).Delete(&models.Campaign{}).Error; err != nil {
			return err
		}

		return nil
	})
}
