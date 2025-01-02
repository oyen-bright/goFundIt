package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

// Constants
const (
	CampaignStatusActive   = "active"
	CampaignStatusEnded    = "ended"
	CampaignStatusUpcoming = "upcoming"
)

// Struct
type Campaign struct {
	ID           string  `gorm:"type:text;primaryKey" validate:"-" binding:"-" json:"id"`
	Key          string  `gorm:"-" validate:"-" binding:"-" json:"key"`
	Title        string  `gorm:"type:varchar(255);not null" encrypt:"true" validate:"required,min=4" binding:"required" json:"title"`
	Description  string  `gorm:"type:text"  encrypt:"true" validate:"required,min=100" binding:"required,min=100" json:"description"`
	TargetAmount float64 `gorm:"not null" validate:"required,gt=0" binding:"-" json:"targetAmount"`

	//Payment
	PaymentMethod PaymentMethod `gorm:"type:varchar(10);not null" validate:"required,oneof=fiat crypto manual" binding:"required,oneof=fiat crypto manual" json:"paymentMethod"`
	FiatCurrency  *FiatCurrency `gorm:"type:varchar(3)" validate:"required_if=PaymentMethod fiat,omitempty" binding:"required_if=PaymentMethod fiat" json:"fiatCurrency,omitempty"`
	CryptoToken   *CryptoToken  `gorm:"type:varchar(10)" validate:"required_if=PaymentMethod crypto,omitempty" binding:"required_if=PaymentMethod crypto" json:"cryptoToken,omitempty"`

	//Relations
	Images       []CampaignImage `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" validate:"-" binding:"omitempty,dive,required" json:"images"`
	Activities   []Activity      `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" validate:"-" binding:"omitempty,dive,required"`
	Contributors []Contributor   `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" validate:"required,gt=0,dive,required" binding:"required,gt=0,dive,required" json:"contributors"`

	//Payout
	Payout *Payout `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" validate:"-" binding:"-" json:"payout"`

	StartDate time.Time `gorm:"not null" validate:"required" binding:"required" json:"startDate"`
	EndDate   time.Time `gorm:"not null" validate:"required,gtfield=StartDate" binding:"required,gtfield=StartDate" json:"endDate"`

	CreatedByHandle string `gorm:"not null" validate:"required" binding:"-" json:"createdByHandle"`
	CreatedBy       User   `gorm:"references:Handle" validate:"-" binding:"-" json:"-"`

	CreatedAt time.Time `gorm:"not null" validate:"-" binding:"-" json:"-"`
	UpdatedAt time.Time `validate:"-" binding:"-" json:"-"`
}

// NewCampaign
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

	for i := range c.Activities {
		c.Activities[i].UpdateCampaignId(c.ID)
		c.Activities[i].ApproveActivity()
		c.Activities[i].UpdateCreatedBy(CreatedBy)
	}

	for i := range c.Contributors {
		c.Contributors[i].UpdateCampaignId(c.ID)
	}

	for i := range c.Images {
		c.Images[i].UpdateCampaignId(c.ID)
	}

	c.CreatedByHandle = CreatedBy.Handle
	c.CreatedBy = CreatedBy
}

// Struct Methods
func (c *Campaign) CanCleanUp() bool {
	hasPayout := c.Payout != nil && c.Payout.Status == PayoutStatusCompleted
	hasPaidContributions := false

	for _, contributor := range c.Contributors {
		if contributor.HasPaid() {
			if !hasPayout {
				return false
			}
			hasPaidContributions = true
			break
		}
	}

	return !hasPaidContributions || hasPayout
}

func (c *Campaign) HasEnded() bool {
	return time.Now().After(c.EndDate)
}

func (c *Campaign) HasStarted() bool {
	return time.Now().After(c.StartDate) || time.Now().Equal(c.StartDate)
}

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

func (c *Campaign) TimeRemaining() time.Duration {
	if c.HasEnded() {
		return 0
	}
	return time.Until(c.EndDate)
}

func (c *Campaign) GetContributorsEmails() []string {

	emails := make([]string, len(c.Contributors))
	for i, contributor := range c.Contributors {
		emails[i] = contributor.Email
	}
	return emails
}

func (c *Campaign) GetContributorByEmail(email string) *Contributor {
	for _, contributor := range c.Contributors {
		if contributor.Email == email {
			return &contributor
		}
	}
	return nil
}

func (c *Campaign) GetContributorByID(ID uint) *Contributor {
	for _, contributor := range c.Contributors {
		if contributor.ID == ID {
			return &contributor
		}
	}
	return nil
}

func (c *Campaign) CanInitiatePayout() bool {
	for _, contributor := range c.Contributors {
		if !contributor.HasPaid() {
			return false
		}
	}
	return true
}

func (c *Campaign) GetPayoutAmount() float64 {
	amount := 0.0
	for _, contributor := range c.Contributors {
		if contributor.HasPaid() {
			amount += contributor.Payment.Amount
		}
	}
	return amount
}

func (c *Campaign) HasReached50PercentMilestone() bool {
	c.UpdateTotalContributionsAmount()
	percentage := (c.GetPayoutAmount() / c.TargetAmount) * 100
	return percentage >= 50
}

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

func (c *Campaign) GetActivityById(ID uint) *Activity {
	for _, activity := range c.Activities {
		if activity.ID == ID {
			return &activity
		}
	}
	return nil
}

func (c *Campaign) Validate() error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil
}

func (c *Campaign) Update(title, description *string, endDate *time.Time) {
	if title != nil {
		c.Title = *title
	}
	if description != nil {
		c.Description = *description
	}
	if endDate != nil {
		c.EndDate = *endDate
	}
}
func (c *Campaign) UpdateTotalContributionsAmount() {
	var totalAmount float64
	for _, contributor := range c.Contributors {
		totalAmount += contributor.GetAmountTotal()
	}
	c.TargetAmount = totalAmount
}

// GORM Hooks ---------------------------------------------------------

func (c *Campaign) BeforeSave(tx *gorm.DB) (err error) {
	c.UpdateTotalContributionsAmount()
	return nil
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	c.UpdateTotalContributionsAmount()

	return c.Validate()
}

// Helper Functions --------------------------------------------------

func generateKey() string {
	return utils.GenerateRandomString("GC-", 8)
}

func generateCampaignId(title string) string {
	return utils.GenerateRandomString(title[:2], 9)
}

// Encryption Methods ----------------------------------------------------

func (c *Campaign) Encrypt(e encryption.Encryptor) error {
	encrypted, err := e.EncryptStruct(c, c.Key)
	if err != nil {
		return err
	}

	if campaign, ok := encrypted.(*Campaign); ok {

		*c = *campaign
	}
	return nil
}

func (c *Campaign) Decrypt(e encryption.Encryptor) error {
	encrypted, err := e.DecryptStruct(c, c.Key)
	if err != nil {
		return err
	}

	if campaign, ok := encrypted.(*Campaign); ok {
		*c = *campaign
	}
	return nil
}
