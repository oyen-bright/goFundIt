package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
)

// TODO: consider required the userEmail via binding to make UpdateUserEmail() redundant
type Contributor struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CampaignID    string    `gorm:"not null;foreignKey:CampaignID;index:idx_campaign_user,unique" validate:"required" json:"campaignId"`
	Amount        float64   `gorm:"not null" binding:"required" validate:"gt=0,required" json:"amount"`
	Email         string    `gorm:"-" binding:"required,email,lowercase" validate:"email,required" json:"email"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"paymentStatus" binding:"-"`
	UserEmail     string    `gorm:"not null;foreignKey:UserEmail;index:idx_campaign_user,unique" json:"userEmail"`
	CreatedAt     time.Time `gorm:"not null" json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (c *Contributor) ToJSON() map[string]interface{} {
	return ToJSON(*c)
}

func (c *Contributor) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

func NewContributor(campaignID, email string, amount float64) *Contributor {
	return &Contributor{
		CampaignID: campaignID,
		Amount:     amount,
		Email:      email,
		UserEmail:  email,
	}
}

// UpdateCampaignId updates the campaignId of the model
//
//   - Note: campaignId is only updated if the current CampaignID is empty
func (c *Contributor) UpdateCampaignId(id string) {

	if c.CampaignID != "" {
		return
	}
	c.CampaignID = id
}

// UpdateUserEmail updates the UserEmail of the model
//
//   - Note: UserEmail is only updated if the current UserEmail is empty
func (c *Contributor) UpdateUserEmail() {
	if c.UserEmail != "" {
		return
	}
	c.UserEmail = c.Email
}

func (c *Contributor) BeforeCreate(tx *gorm.DB) (err error) {

	if validationErrors := c.Validate(); validationErrors != nil {
		return validationErrors
	}
	return nil
}
