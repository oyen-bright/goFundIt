package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type campaignService struct {
	repo               repositories.CampaignRepository
	contributorService services.ContributorService
	authService        services.AuthService
	logger             logger.Logger
}

func NewCampaignService(repo repositories.CampaignRepository, contributorService services.ContributorService, authService services.AuthService, logger logger.Logger) services.CampaignService {
	return &campaignService{repo: repo, contributorService: contributorService, authService: authService, logger: logger}
}

// CreateCampaign creates a new campaign for a user.
func (s *campaignService) CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error) {

	// check if user can create a Campaign via the user handle
	if err := s.UserCanCreateCampaign(userHandle); err != nil {
		return models.Campaign{}, err
	}

	// Check if contributors is not involved in any other campaigns
	emailsCanNotContribute, err := s.EmailsCanContribute(campaign.GetContributorsEmails())
	if err != nil {
		return models.Campaign{}, err
	}
	if len(emailsCanNotContribute) > 0 {
		return models.Campaign{}, errs.BadRequest("Emails "+strings.Join(emailsCanNotContribute, ", ")+" already have campaigns.", emailsCanNotContribute)
	}

	// Get the user by the user handle
	user, err := s.authService.GetUserByHandle(userHandle)
	if err != nil {
		return models.Campaign{}, err
	}

	// Create CampaignID and update the CampaignID on the bound images, activities, and contributors as well as update the creeatedBy of the campaign and activities
	campaign.FromBinding(user)

	// Create new users struct from the contributors
	campaignUsers := getUsersFromContributors(campaign.Contributors)

	// Remove already existing users
	newUsers, err := s.authService.FindNonExistingUsers(campaignUsers)
	if err != nil {
		return models.Campaign{}, err
	}

	// Create the new users
	_, err = s.authService.CreateUsers(newUsers)
	if err != nil {
		return models.Campaign{}, err
	}

	// Create campaign
	campaign, err = s.repo.Create(&campaign)
	if err != nil {
		return models.Campaign{}, errs.InternalServerError(err).Log(s.logger)
	}

	return campaign, nil
}

// GetCampaignByID fetches campaign by id is not found returns err
func (s *campaignService) GetCampaignByID(id string) (*models.Campaign, error) {
	campaign, err := s.repo.GetByID(id, true)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil, errs.BadRequest("Campaign does not exist", nil)
		}
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	fmt.Println(campaign.CreatedBy.Email)
	return &campaign, nil
}

// userCanCreateCampaign checks if a user is allowed to create a campaign.
func (s *campaignService) UserCanCreateCampaign(userHandle string) error {

	_, err := s.repo.GetByCreatorHandle(userHandle, false)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil
		}

		return errs.InternalServerError(err).Log(s.logger)
	}
	return errs.BadRequest("You have an existing campaign. Please finish it before creating a new one.", nil)
}

// emailsCanContribute checks if the provided email addresses can contribute to a campaign.
func (s *campaignService) EmailsCanContribute(contributorsEmail []string) ([]string, error) {
	return s.contributorService.GetEmailsOfExistingContributors(contributorsEmail)
}

// getUsersFromContributors converts a list of contributors to a list of users.
func getUsersFromContributors(contributors []models.Contributor) []models.User {
	users := make([]models.User, 0, len(contributors))

	for _, c := range contributors {
		users = append(users, *models.NewUser("", c.Email, false))
	}

	return users
}

func validateContributionSum(fl validator.FieldLevel) bool {
	campaign, ok := fl.Parent().Interface().(models.Campaign)
	if !ok {
		return false
	}

	totalContributions := calculateContributionTotal(campaign)

	return totalContributions == campaign.TargetAmount
}

func calculateContributionTotal(campaign models.Campaign) float64 {
	var totalAmount float64
	for _, contributor := range campaign.Contributors {
		totalAmount += contributor.Amount
	}

	return totalAmount
}

func isCampaignStartDateValid(campaign models.Campaign) bool {
	return campaign.StartDate.After(time.Now()) || campaign.StartDate.Equal(time.Now())
}
