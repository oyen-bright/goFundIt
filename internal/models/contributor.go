package models

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Payment status constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
)

// Contributor represents a user who contributes funds to a campaign
// TODO: use DTO to bind the email and name
// TODO: consider required the userEmail via binding to make UpdateUserEmail() redundant
type Contributor struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `gorm:"type:varchar(255);default:null;null" validate:"omitempty" binding:"omitempty,gte=3" json:"name"`
	CampaignID string  `gorm:"not null;foreignKey:CampaignID;index:idx_campaign_user,unique" validate:"required" json:"campaignId"`
	Amount     float64 `gorm:"not null" binding:"required" validate:"gt=0,required" json:"amount"`
	// Email         string    `gorm:"-" binding:"required,email,lowercase" json:"email"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"paymentStatus" binding:"-"`
	Email         string    `gorm:"not null;foreignKey:Email;index:idx_campaign_user,unique" json:"email" binding:"-"`
	CreatedAt     time.Time `gorm:"not null" json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// Constructor

// NewContributor creates a new Contributor instance with the provided parameters
func NewContributor(campaignID, email string, amount float64) *Contributor {
	return &Contributor{
		CampaignID: campaignID,
		Amount:     amount,
		Email:      email,
		// UserEmail:  email,
	}
}

// Payment Status Methods

// HasPaid checks if the payment has been successfully processed
func (c *Contributor) HasPaid() bool {
	return c.PaymentStatus == PaymentStatusSucceeded
}

// IsPending checks if the payment is still pending
func (c *Contributor) IsPending() bool {
	return c.PaymentStatus == PaymentStatusPending
}

// HasFailed checks if the payment has failed
func (c *Contributor) HasFailed() bool {
	return c.PaymentStatus == PaymentStatusFailed
}

// Update Methods

// UpdateCampaignId updates the campaignId of the model
// Note: campaignId is only updated if the current CampaignID is empty
func (c *Contributor) UpdateCampaignId(id string) {
	if c.CampaignID != "" {
		return
	}
	c.CampaignID = id
}

// UpdateUserEmail updates the UserEmail of the model
// Note: UserEmail is only updated if the current UserEmail is empty
func (c *Contributor) UpdateUserEmail() {
	// if c.UserEmail != "" {
	// 	return
	// }
	// c.UserEmail = c.Email

}

// Validation Methods

// Validate performs validation checks on the contributor
func (c *Contributor) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

// GORM Hooks

// BeforeCreate performs validation before creating the contributor
func (c *Contributor) BeforeCreate(tx *gorm.DB) (err error) {
	if validationErrors := c.Validate(); validationErrors != nil {
		return validationErrors
	}
	return nil
}

// BeforeSave GORM hook to convert empty string to NULL
func (c *Contributor) BeforeSave(tx *gorm.DB) error {
	if c.Name == "" {
		c.Name = sql.NullString{}.String
	}
	return nil
}
