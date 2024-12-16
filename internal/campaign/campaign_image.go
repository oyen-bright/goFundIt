package campaign

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CampaignImage struct {
	ID         uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	CampaignID string `gorm:"not null;foreignKey:CampaignID" validate:"required" json:"-"`
	ImageUrl   string `gorm:"type:varchar(255);not null" validate:"required,url" binding:"required,url" json:"imageUrl"`
}

func (c *CampaignImage) ToJSON() string {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func NewImage(campaignID, imageURL string) *CampaignImage {
	return &CampaignImage{
		ImageUrl:   imageURL,
		CampaignID: campaignID,
	}
}

func (c *CampaignImage) Validate() error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil

}

// UpdateCampaignId updates the campaignId of the model
//
//   - Note: campaignId is only updated if the currentCampaignID is empty
func (a *CampaignImage) UpdateCampaignId(id string) {

	if a.CampaignID != "" {
		return
	}
	a.CampaignID = id
}

func (c *CampaignImage) BeforeCreate(tx *gorm.DB) (err error) {
	return c.Validate()
}
