package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

// TODO: consider using DTOs for the models
type Campaign struct {
	ID              string          `gorm:"type:text;primaryKey" json:"id"`
	Key             string          `gorm:"-" json:"key"`
	Title           string          `gorm:"type:varchar(255);not null" validate:"required,min=4" binding:"required" json:"title"`
	Description     string          `gorm:"type:text" binding:"required" validate:"required,min=100" json:"description"`
	TargetAmount    float64         `gorm:"not null" validate:"required,gt=0" binding:"required,gt=0" json:"targetAmount"`
	Images          []CampaignImage `gorm:"foreignKey:CampaignID" binding:"omitempty,dive,required" validate:"-" json:"images"`
	Activities      []Activity      `gorm:"foreignKey:CampaignID" binding:"omitempty,dive,required" validate:"-" json:"activities"`
	Contributors    []Contributor   `gorm:"many2many:campaign_contributors" binding:"required,gt=0,dive,required" validate:"required,gt=0,dive,required,contributorSum" json:"contributors"`
	StartDate       time.Time       `gorm:"not null" validate:"required" binding:"required" json:"startDate"`
	EndDate         time.Time       `gorm:"not null" validate:"required,gtfield=StartDate" binding:"required,gtfield=StartDate" json:"endDate"`
	CreatedByHandle string          `gorm:"not null" validate:"required" binding:"-" json:"createdByHandle"`
	CreatedBy       User            `gorm:"references:Handle" binding:"-" validate:"-" json:"-"`
	CreatedAt       time.Time       `gorm:"not null" json:"-"`
	UpdatedAt       time.Time       `json:"-"`
}

func NewCampaign(title, description string, targetAmount float64, startDate, endDate time.Time, images []CampaignImage, activities []Activity, contributors []Contributor, CreatedBy User) *Campaign {
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

// func (c *Campaign) ToJSON() map[string]interface{} {
// 	return ToJSON(*c)
// }

// FromBinding initializes a Campaign instance with the provided auth.User as the creator.
// It generates a unique campaign ID based on the campaign's title and updates all associated
// activities, contributors, and images with this ID. Additionally, it approves all activities
// and sets the creator's handle and user information in the campaign.
//
// Parameters:
//   - CreatedBy: The auth.User who created the campaign.
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

func (c *Campaign) EmailIsPartOfCampaign(email string) bool {

	if c.CreatedBy.Email == email {
		return true
	}

	for _, contributor := range c.Contributors {

		fmt.Println(contributor.UserEmail)
		if contributor.UserEmail == email {
			return true
		}
	}
	return false
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

func validateContributionSum(fl validator.FieldLevel) bool {
	campaign, ok := fl.Parent().Interface().(Campaign)
	if !ok {
		return false
	}

	totalContributions := calculateContributionTotal(campaign)

	return totalContributions == campaign.TargetAmount
}

func calculateContributionTotal(campaign Campaign) float64 {
	var totalAmount float64
	for _, contributor := range campaign.Contributors {
		totalAmount += contributor.Amount
	}

	return totalAmount
}

func isCampaignStartDateValid(campaign Campaign) bool {
	return campaign.StartDate.After(time.Now()) || campaign.StartDate.Equal(time.Now())
}

func NewDummyCampaign() *Campaign {
	images := []CampaignImage{
		{ImageUrl: "https://example.com/image1.jpg"},
		{ImageUrl: "https://example.com/image2.jpg"},
	}

	activities := []Activity{
		{Title: "Activity 1", Subtitle: "Description for activity 1"},
		{Title: "Activity 2", Subtitle: "Description for activity 2"},
	}

	contributors := []Contributor{
		{Email: "contributor1@example.com", Amount: 50.0},
		{Email: "contributor2@example.com", Amount: 100.0},
	}

	createdBy := User{
		Handle: "user123",
		Email:  "user123@example.com",
	}

	return NewCampaign(
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
