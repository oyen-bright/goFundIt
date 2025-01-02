package models

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TODO: use DTO to bind the email and name
// TODO: add notification settings - allow contributor to opt-in for notifications or not

// Original struct for reference
// type Contributor struct {
// 	ID            uint          `gorm:"primaryKey" json:"id"`
// 	Name          string        `gorm:"type:varchar(255);default:null;null" validate:"omitempty" binding:"omitempty,gte=3" json:"name"`
// 	CampaignID    string        `gorm:"not null;foreignKey:CampaignID;index:idx_campaign_user,unique" validate:"required" json:"campaignId"`
// 	Amount        float64       `gorm:"not null" binding:"required,gte=0" validate:"gte=0,required" json:"amount"`
// 	AmountPaid    *float64      `gorm:"type:decimal(10,2);default:null" binding:"-" json:"amountPaid,omitempty"`
// 	Activities    []Activity    `gorm:"many2many:activities_contributors" binding:"-" json:"activities"`
// 	PaymentStatus PaymentStatus `gorm:"not null;default:pending" json:"paymentStatus" binding:"-"`
// 	Email         string        `gorm:"not null;foreignKey:Email;index:idx_campaign_user,unique" json:"email" binding:"-"`
// 	CreatedAt     time.Time     `gorm:"not null" json:"-"`
// 	UpdatedAt     time.Time     `json:"-"`
// }

type Contributor struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `gorm:"type:varchar(255);default:null;null" validate:"omitempty" binding:"omitempty,gte=3" json:"name"`
	CampaignID string     `gorm:"not null;foreignKey:CampaignID;index:idx_campaign_user,unique" validate:"required" json:"campaignId"`
	Amount     float64    `gorm:"not null" binding:"required,gte=0" validate:"gte=0,required" json:"amount"`
	Activities []Activity `gorm:"many2many:activities_contributors" binding:"-" json:"activities"`

	Payment *Payment `gorm:"foreignKey:ContributorID" json:"payment"`

	Email     string    `gorm:"not null;foreignKey:Email;index:idx_campaign_user,unique" json:"email" binding:"-"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Constructor
func NewContributor(campaignID, email string, amount float64) *Contributor {
	return &Contributor{
		CampaignID: campaignID,
		Amount:     amount,
		Email:      email,
		// UserEmail:  email,
	}
}

// Payment Status Methods
func (c *Contributor) HasPaid() bool {
	if c.Payment == nil {
		return false
	}
	//TODO: not accountable for all cases
	return c.Payment.PaymentStatus == PaymentStatusSucceeded || c.Payment.PaymentStatus == PaymentStatusPendingApproval
}

func (c *Contributor) IsPending() bool {
	if c.Payment == nil {
		return false
	}
	return c.Payment.PaymentStatus == PaymentStatusPending
}

func (c *Contributor) HasFailed() bool {
	if c.Payment == nil {
		return false
	}
	return c.Payment.PaymentStatus == PaymentStatusFailed
}

// Amount Methods
func (c *Contributor) GetAmountTotal() float64 {
	for _, activity := range c.Activities {
		c.Amount += activity.Cost
	}
	return c.Amount
}

// Update Methods
func (c *Contributor) UpdateCampaignId(id string) {
	if c.CampaignID != "" {
		return
	}
	c.CampaignID = id
}

// Validation Methods
func (c *Contributor) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

// GORM Hooks
func (c *Contributor) BeforeCreate(tx *gorm.DB) (err error) {
	if validationErrors := c.Validate(); validationErrors != nil {
		return validationErrors
	}
	return nil
}

func (c *Contributor) BeforeSave(tx *gorm.DB) error {
	if c.Name == "" {
		c.Name = sql.NullString{}.String
	}
	return nil
}
