package campaign

import (
	"strings"

	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/contributor"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type CampaignService interface {
	//  creates a new campaign for a user.
	CreateCampaign(campaign Campaign, userHandle string) (Campaign, error)

	userCanCreateCampaign(userHandle string) error
	emailsCanContribute(contributorsEmail []string) ([]string, error)
	// GetCampaignByID(id string) (Campaign, error)
	// UpdateCampaign(campaign Campaign) (Campaign, error)
	// DeleteCampaign(campaignID string) error
	// JoinCampaign(userID uint, campaignID string) error
	// LeaveCampaign(userID uint, campaignID string) error
	// SetCampaignTheme(campaignID string, themeID uint) error
	// GetCampaignsByUser(handle string) ([]Campaign, error)
	// UserCanContribute(handle string, campaignID string) error
}

type campaignService struct {
	repo               CampaignRepository
	contributorService contributor.ContributorService
	authService        auth.AuthService
	logger             logger.Logger
}

func Service(repo CampaignRepository, contributorService contributor.ContributorService, authService auth.AuthService, logger logger.Logger) CampaignService {
	return &campaignService{repo: repo, contributorService: contributorService, authService: authService, logger: logger}
}

// CreateCampaign creates a new campaign for a user.
func (s *campaignService) CreateCampaign(campaign Campaign, userHandle string) (Campaign, error) {

	// check if user can create a Campaign via the user handle
	if err := s.userCanCreateCampaign(userHandle); err != nil {
		return Campaign{}, err
	}

	// Check if contributors is not involved in any other campaigns
	emailsCanNotContribute, err := s.emailsCanContribute(campaign.GetContributorsEmails())
	if err != nil {
		return Campaign{}, err
	}
	if len(emailsCanNotContribute) > 0 {
		return Campaign{}, errs.BadRequest("Emails "+strings.Join(emailsCanNotContribute, ", ")+" already have campaigns.", emailsCanNotContribute)
	}

	// Get the user by the user handle
	user, err := s.authService.GetUserByHandle(userHandle)
	if err != nil {
		return Campaign{}, err
	}

	// Create CampaignID and update the CampaignID on the bound images, activities, and contributors as well as update the creeatedBy of the campaign and activities
	campaign.FromBinding(user)

	// Create new users struct from the contributors
	campaignUsers := getUsersFromContributors(campaign.Contributors)

	// Remove already existing users
	newUsers, err := s.authService.FindNonExistingUsers(campaignUsers)
	if err != nil {
		return Campaign{}, err
	}

	// Create the new users
	_, err = s.authService.CreateUsers(newUsers)
	if err != nil {
		return Campaign{}, err
	}

	// Create campaign
	campaign, err = s.repo.Create(&campaign)
	if err != nil {
		return Campaign{}, errs.InternalServerError(err).Log(s.logger)
	}

	return campaign, nil

}

// userCanCreateCampaign checks if a user is allowed to create a campaign.
func (s *campaignService) userCanCreateCampaign(userHandle string) error {

	_, err := s.repo.GetByCreatorHandle(userHandle, false)
	if err != nil {
		if errs.NewDB(err).IsNotfound() {
			return nil
		}

		return errs.InternalServerError(err).Log(s.logger)
	}
	return errs.BadRequest("You have an existing campaign. Please finish it before creating a new one.", nil)
}

// emailsCanContribute checks if the provided email addresses can contribute to a campaign.
func (s *campaignService) emailsCanContribute(contributorsEmail []string) ([]string, error) {
	return s.contributorService.GetEmailsOfExistingContributors(contributorsEmail)
}

// getUsersFromContributors converts a list of contributors to a list of users.
func getUsersFromContributors(contributors []contributor.Contributor) []auth.User {
	users := make([]auth.User, 0, len(contributors))

	for _, c := range contributors {
		users = append(users, *auth.New("", c.Email, false))
	}

	return users
}
