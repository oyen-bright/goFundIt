package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Activity struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	CampaignID      string        `gorm:"type:text;foreignKey:CampaignID" validate:"required" json:"campaignId"`
	Title           string        `gorm:"type:varchar(255);not null" binding:"required,min=4" json:"title"`
	Subtitle        string        `gorm:"type:varchar(255)" json:"subtitle"`
	ImageUrl        string        `gorm:"type:varchar(255)" binding:"omitempty,url" validate:"omitempty,url" json:"imageUrl"`
	IsMandatory     bool          `gorm:"not null" binding:"boolean" json:"isMandatory"`
	Cost            float64       `gorm:"not null" binding:"required" validate:"required,gt=0" json:"cost"`
	IsApproved      bool          `gorm:"not null; default:false" json:"isApproved"`
	Contributors    []Contributor `gorm:"many2many:activity_contributors" binding:"-" json:"contributors"`
	CreatedByHandle string        `gorm:"not null" validate:"required" json:"-"`
	CreatedBy       User          `gorm:"references:Handle"  binding:"-" validate:"-" json:"-"`
	CreatedAt       time.Time     `gorm:"not null" json:"-"`
	UpdatedAt       time.Time     `json:"-"`
}

func (a *Activity) ToJSON() map[string]interface{} {
	return ToJSON(*a)
}

func (a *Activity) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(a); err != nil {
		return err
	}
	return nil
}

func New(campaignID, title, subtitle, imageURL, CreatedByHandle string, isMandatory, isApproved bool, cost float64) *Activity {
	return &Activity{
		CampaignID:      campaignID,
		Title:           title,
		Subtitle:        subtitle,
		ImageUrl:        imageURL,
		Cost:            cost,
		IsMandatory:     isMandatory,
		IsApproved:      isApproved,
		CreatedByHandle: CreatedByHandle,
	}
}

func (a *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	return a.Validate()
}

// UpdateCampaignId updates the campaignId of the model
//
//   - Note: campaignId is only updated if the currentCampaignID is empty
func (a *Activity) UpdateCampaignId(id string) {

	if a.CampaignID != "" {
		return
	}
	a.CampaignID = id
}

// UpdateCreatedByHandle updates the CreatedByHandle of the model
//
//   - Note: CreatedByHandle is only updated if the current CreatedByHandle is empty
func (a *Activity) UpdateCreatedBy(user User) {

	if a.CreatedByHandle != "" {
		return
	}
	a.CreatedByHandle = user.Handle
	a.CreatedBy = user

}

// ApproveActivity sets the IsApproved field to true
func (a *Activity) ApproveActivity() {
	a.IsApproved = true
}

// ApproveActivity sets the IsApproved field to true
func (a *Activity) MarkAsNotMandatory() {
	a.IsMandatory = false
}
