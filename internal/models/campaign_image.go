package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Struct definition
type CampaignImage struct {
	ID         uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	CampaignID string `gorm:"not null;foreignKey:CampaignID" validate:"required" json:"-"`
	ImageUrl   string `gorm:"type:varchar(255);not null" validate:"required,url" binding:"required,url" json:"imageUrl"`
}

// Constructor

func NewImage(campaignID, imageURL string) *CampaignImage {
	return &CampaignImage{
		ImageUrl:   imageURL,
		CampaignID: campaignID,
	}
}

// Methods

func (c *CampaignImage) ToJSON() map[string]interface{} {
	return ToJSON(*c)
}

func (c *CampaignImage) Validate() error {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

func (c *CampaignImage) UpdateCampaignId(id string) {
	if c.CampaignID != "" {
		return
	}
	c.CampaignID = id
}

// GORM Hooks

func (c *CampaignImage) BeforeCreate(tx *gorm.DB) (err error) {
	return c.Validate()
}
