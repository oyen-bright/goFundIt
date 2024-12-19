package models

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

// Campaign status constants
const (
	CampaignStatusActive   = "active"
	CampaignStatusEnded    = "ended"
	CampaignStatusUpcoming = "upcoming"
)

// Campaign represents a fundraising campaign
type Campaign struct {
	ID              string          `gorm:"type:text;primaryKey" json:"id"`
	Key             string          `gorm:"-" json:"key"`
	Title           string          `gorm:"type:varchar(255);not null" validate:"required,min=4" binding:"required" json:"title"`
	Description     string          `gorm:"type:text" binding:"required" validate:"required,min=100" json:"description"`
	TargetAmount    float64         `gorm:"not null" validate:"required,gt=0" binding:"required,gt=0" json:"targetAmount"`
	Images          []CampaignImage `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" binding:"omitempty,dive,required" validate:"-" json:"images"`
	Activities      []Activity      `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" binding:"omitempty,dive,required" validate:"-" `
	Contributors    []Contributor   `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" binding:"required,gt=0,dive,required" validate:"required,gt=0,dive,required,contributorSum" json:"contributors"`
	StartDate       time.Time       `gorm:"not null" validate:"required" binding:"required" json:"startDate"`
	EndDate         time.Time       `gorm:"not null" validate:"required,gtfield=StartDate" binding:"required,gtfield=StartDate" json:"endDate"`
	CreatedByHandle string          `gorm:"not null" validate:"required" binding:"-" json:"createdByHandle"`
	CreatedBy       User            `gorm:"references:Handle" binding:"-" validate:"-" json:"-"`
	CreatedAt       time.Time       `gorm:"not null" json:"-"`
	UpdatedAt       time.Time       `json:"-"`
}

// Campaign Status Methods

// HasEnded checks if the campaign has ended based on current time
func (c *Campaign) HasEnded() bool {
	return time.Now().After(c.EndDate)
}

// HasStarted checks if the campaign has started based on current time
func (c *Campaign) HasStarted() bool {
	return time.Now().After(c.StartDate) || time.Now().Equal(c.StartDate)
}

// GetStatus returns the current status of the campaign
func (c *Campaign) GetStatus() string {
	now := time.Now()

	if now.Before(c.StartDate) {
		return CampaignStatusUpcoming
	}

	if now.After(c.EndDate) {
		return CampaignStatusEnded
	}

	return CampaignStatusActive
}

// TimeRemaining returns the duration until campaign ends
// Returns 0 if campaign has ended
func (c *Campaign) TimeRemaining() time.Duration {
	if c.HasEnded() {
		return 0
	}
	return time.Until(c.EndDate)
}

// Contributor Related Methods

// GetContributorsEmails returns a slice of all contributor emails
func (c *Campaign) GetContributorsEmails() []string {
	emails := make([]string, len(c.Contributors))
	for i, contributor := range c.Contributors {
		emails[i] = contributor.Email
	}
	return emails
}

// GetContributorByEmail returns a contributor by their email
func (c *Campaign) GetContributorByEmail(email string) *Contributor {
	for _, contributor := range c.Contributors {
		if contributor.Email == email {
			return &contributor
		}
	}
	return nil
}

// GetContributor returns a contributor by their Id
func (c *Campaign) GetContributorByID(ID uint) *Contributor {
	for _, contributor := range c.Contributors {
		if contributor.ID == ID {
			return &contributor
		}
	}
	return nil
}

// EmailIsPartOfCampaign checks if an email is associated with the campaign
func (c *Campaign) EmailIsPartOfCampaign(email string) bool {
	if c.CreatedBy.Email == email {
		return true
	}

	for _, contributor := range c.Contributors {
		if contributor.Email == email {
			return true
		}
	}
	return false
}

// Activity Related Methods

// GetActivityById returns an activity by their id return nil if not found
func (c *Campaign) GetActivityById(ID uint) *Activity {

	for _, activity := range c.Activities {
		if activity.ID == ID {
			return &activity
		}
	}
	return nil
}

// Validate performs all validation checks on the campaign
func (c *Campaign) Validate() error {
	v := validator.New()
	v.RegisterValidation("contributorSum", validateContributionSum)

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil
}

// ValidateNewCampaign performs validation specific to new campaigns
func (c *Campaign) ValidateNewCampaign() error {

	if time.Now().After(c.StartDate) {
		return errors.New("campaign start date must be today or in the future")
	}

	total := calculateContributionTotal(*c)
	if total != c.TargetAmount {
		return errors.New("total contributions do not match the target amount")
	}

	return nil
}

// BeforeCreate GORM hook for validation before creating
func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	return c.Validate()
}

// Constructor and Initialization Methods

// NewCampaign creates a new Campaign instance
func NewCampaign(title, description string, targetAmount float64, startDate, endDate time.Time,
	images []CampaignImage, activities []Activity, contributors []Contributor, CreatedBy User) *Campaign {
	return &Campaign{
		ID:              generateCampaignId(title),
		Key:             generateKey(),
		Title:           title,
		Description:     description,
		Images:          images,
		Activities:      activities,
		Contributors:    contributors,
		TargetAmount:    targetAmount,
		StartDate:       startDate,
		EndDate:         endDate,
		CreatedByHandle: CreatedBy.Handle,
		CreatedBy:       CreatedBy,
	}
}

// FromBinding initializes a Campaign instance from binding data
func (c *Campaign) FromBinding(CreatedBy User) {
	c.ID = generateCampaignId(c.Title)
	c.Key = generateKey()

	// Initialize activities
	for i := range c.Activities {
		c.Activities[i].UpdateCampaignId(c.ID)
		c.Activities[i].ApproveActivity()
		c.Activities[i].UpdateCreatedBy(CreatedBy)
	}

	// Initialize contributors
	for i := range c.Contributors {
		c.Contributors[i].UpdateCampaignId(c.ID)
		c.Contributors[i].UpdateUserEmail()
	}

	// Initialize images
	for i := range c.Images {
		c.Images[i].UpdateCampaignId(c.ID)
	}

	c.CreatedByHandle = CreatedBy.Handle
	c.CreatedBy = CreatedBy
}

// Helper Functions

func generateKey() string {
	return utils.GenerateRandomString("GC-", 8)
}

func generateCampaignId(title string) string {
	return utils.GenerateRandomString(title[:2], 9)
}

func validateContributionSum(fl validator.FieldLevel) bool {
	campaign, ok := fl.Parent().Interface().(Campaign)
	if !ok {
		return false
	}
	return calculateContributionTotal(campaign) == campaign.TargetAmount
}

func calculateContributionTotal(campaign Campaign) float64 {
	var totalAmount float64
	for _, contributor := range campaign.Contributors {
		totalAmount += contributor.Amount
	}
	return totalAmount
}
