package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Activity represents an action or task within a campaign that requires funding
type Activity struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	CampaignID      string        `gorm:"type:text;foreignKey:CampaignID" validate:"required" json:"campaignId"`
	Title           string        `gorm:"type:varchar(255);not null" binding:"required,min=4" json:"title"`
	Subtitle        string        `gorm:"type:varchar(255)" json:"subtitle"`
	ImageUrl        string        `gorm:"type:varchar(255)" binding:"omitempty,url" validate:"omitempty,url" json:"imageUrl"`
	IsMandatory     bool          `gorm:"not null" binding:"boolean" json:"isMandatory"`
	Cost            float64       `gorm:"not null" binding:"required" validate:"required,gt=0" json:"cost"`
	IsApproved      bool          `gorm:"not null; default:false" json:"isApproved"`
	Contributors    []Contributor `gorm:"many2many:activities_contributors" binding:"-" json:"contributors"`
	Comments        []Comment     `gorm:"foreignKey:ActivityID" json:"comments" binding:"-"`
	CreatedByHandle string        `gorm:"not null" validate:"required" json:"-"`
	CreatedBy       User          `gorm:"references:Handle" binding:"-" validate:"-" json:"-"`
	CreatedAt       time.Time     `gorm:"not null" json:"-"`
	UpdatedAt       time.Time     `json:"-"`
}

// Constructor

// New creates a new Activity instance with the provided parameters
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

// Contributor Methods

// GetPaidContributors returns a slice of contributors who have successfully paid
func (a *Activity) GetPaidContributors() []Contributor {
	paidContributors := make([]Contributor, 0)
	for _, contributor := range a.Contributors {
		if contributor.HasPaid() {
			paidContributors = append(paidContributors, contributor)
		}
	}
	return paidContributors
}

// GetPaidContributorsCount returns the number of paid contributors
func (a *Activity) GetPaidContributorsCount() int {
	return len(a.GetPaidContributors())
}

// Update Methods

// UpdateCampaignId updates the campaignId of the model
// Note: campaignId is only updated if the currentCampaignID is empty
func (a *Activity) UpdateCampaignId(id string) {
	if a.CampaignID != "" {
		return
	}
	a.CampaignID = id
}

// UpdateCreatedBy updates the CreatedByHandle and CreatedBy of the model
// Note: CreatedByHandle is only updated if the current CreatedByHandle is empty
func (a *Activity) UpdateCreatedBy(user User) {
	if a.CreatedByHandle != "" {
		return
	}
	a.CreatedByHandle = user.Handle
	a.CreatedBy = user
}

// Activity Methods

// ApproveActivity sets the IsApproved field to true
func (a *Activity) ApproveActivity() {
	a.IsApproved = true
}

// MarkAsNotMandatory sets the IsMandatory field to false
func (a *Activity) MarkAsNotMandatory() {
	a.IsMandatory = false
}

// AddContributor adds a contributor to the activity
func (a *Activity) AddContributor(contributor Contributor) {
	a.Contributors = append(a.Contributors, contributor)
}

// RemoveContributor removes a contributor from the activity
func (a *Activity) RemoveContributor(contributor Contributor) {
	for i, c := range a.Contributors {
		if c.ID == contributor.ID {
			a.Contributors = append(a.Contributors[:i], a.Contributors[i+1:]...)
			break
		}
	}
}

// IsContributorOptedIn checks if a contributor has opted in to the activity
func (a *Activity) IsContributorOptedIn(contributorID uint) bool {
	for _, contributor := range a.Contributors {
		if contributor.ID == contributorID {
			return true
		}
	}
	return false
}

// Validation Methods

// Validate performs validation checks on the activity
func (a *Activity) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(a); err != nil {
		return err
	}
	return nil
}

// GORM Hooks

// BeforeCreate performs validation before creating the activity
func (a *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	return a.Validate()
}
