package campaign

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oyen-bright/goFundIt/internal/activity"
	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"github.com/oyen-bright/goFundIt/internal/utils"
	"gorm.io/gorm"
)

type Campaign struct {
	ID              string                    `gorm:"type:text;primaryKey" json:"id"`
	Key             string                    `gorm:"-" json:"key"`
	Title           string                    `gorm:"type:varchar(255);not null" validate:"required,min=4" binding:"required" json:"title"`
	Description     string                    `gorm:"type:text" binding:"required" validate:"required,min=100" json:"description"`
	TargetAmount    float64                   `gorm:"not null" validate:"required,gt=0" binding:"required,gt=0" json:"targetAmount"`
	Images          []CampaignImage           `gorm:"foreignKey:CampaignID" binding:"omitempty,dive,required" validate:"-" json:"images"`
	Activities      []activity.Activity       `gorm:"foreignKey:CampaignID" binding:"omitempty,dive,required" validate:"-" json:"activities"`
	Contributors    []contributor.Contributor `gorm:"many2many:campaign_contributors" binding:"required,gt=0,dive,required" validate:"required,gt=0,dive,required,contributorSum" json:"contributors"`
	StartDate       time.Time                 `gorm:"not null" validate:"required" binding:"required" json:"startDate"`
	EndDate         time.Time                 `gorm:"not null" validate:"required,gtfield=StartDate" binding:"required,gtfield=StartDate" json:"endDate"`
	CreatedByHandle string                    `gorm:"not null" validate:"required" binding:"-" json:"createdByHandle"`
	CreatedBy       auth.User                 `gorm:"references:Handle" binding:"-" validate:"-" json:"-"`
	CreatedAt       time.Time                 `gorm:"not null" json:"-"`
	UpdatedAt       time.Time                 `json:"-"`
}

func New(title, description string, targetAmount float64, startDate, endDate time.Time, images []CampaignImage, activities []activity.Activity, contributors []contributor.Contributor, CreatedBy auth.User) *Campaign {
	c := &(Campaign{
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
	})
	return c
}

func (c *Campaign) ToJSON() string {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// FromBinding initializes a Campaign instance with the provided auth.User as the creator.
// It generates a unique campaign ID based on the campaign's title and updates all associated
// activities, contributors, and images with this ID. Additionally, it approves all activities
// and sets the creator's handle and user information in the campaign.
//
// Parameters:
//   - CreatedBy: The auth.User who created the campaign.
func (c *Campaign) FromBinding(CreatedBy auth.User) {
	c.ID = generateCampaignId(c.Title)
	c.Key = generateKey()
	for i := range c.Activities {
		c.Activities[i].UpdateCampaignId(c.ID)
		c.Activities[i].ApproveActivity()
		c.Activities[i].UpdateCreatedBy(CreatedBy)

	}

	for i := range c.Contributors {
		c.Contributors[i].UpdateCampaignId(c.ID)
		c.Contributors[i].UpdateUserEmail()
	}

	for i := range c.Images {
		c.Images[i].UpdateCampaignId(c.ID)
	}

	c.CreatedByHandle = CreatedBy.Handle
	c.CreatedBy = CreatedBy
}

func (c *Campaign) GetContributorsEmails() []string {
	emails := make([]string, len(c.Contributors))
	for i, contributor := range c.Contributors {
		emails[i] = contributor.Email
	}
	return emails
}

func (c *Campaign) Validate() error {
	v := validator.New()
	v.RegisterValidation("contributorSum", validateContributionSum)

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil
}

func (c *Campaign) ValidateNewCampaign() error {

	isValid := isCampaignStartDateValid(*c)
	if !isValid {
		return errors.New("campaign start date is invalid")
	}
	total := calculateContributionTotal(*c)

	isValid = total == c.TargetAmount

	if !isValid {
		return errors.New("total contributions do not match the target amount")
	}

	return nil
}

func generateKey() string {
	return utils.GenerateRandomString("GC-", 8)
}

func generateCampaignId(title string) string {
	return utils.GenerateRandomString(title[:2], 9)
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	return c.Validate()
}

func NewDummyCampaign() *Campaign {
	images := []CampaignImage{
		{ImageUrl: "https://example.com/image1.jpg"},
		{ImageUrl: "https://example.com/image2.jpg"},
	}

	activities := []activity.Activity{
		{Title: "Activity 1", Subtitle: "Description for activity 1"},
		{Title: "Activity 2", Subtitle: "Description for activity 2"},
	}

	contributors := []contributor.Contributor{
		{Email: "contributor1@example.com", Amount: 50.0},
		{Email: "contributor2@example.com", Amount: 100.0},
	}

	createdBy := auth.User{
		Handle: "user123",
		Email:  "user123@example.com",
	}

	return New(
		"Dummy Campaign",
		"This is a dummy campaign for testing purposes.",
		150.0,
		time.Now(),
		time.Now().AddDate(0, 1, 0),
		images,
		activities,
		contributors,
		createdBy,
	)
}
