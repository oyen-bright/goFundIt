package services

import (
	"fmt"
	"strings"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type campaignService struct {
	repo                repositories.CampaignRepository
	authService         services.AuthService
	notificationService services.NotificationService
	broadcaster         services.EventBroadcaster
	logger              logger.Logger
}

func NewCampaignService(repo repositories.CampaignRepository, authService services.AuthService, notificationService services.NotificationService, broadcast services.EventBroadcaster, logger logger.Logger) services.CampaignService {
	return &campaignService{repo: repo, authService: authService, logger: logger, broadcaster: broadcast, notificationService: notificationService}
}

// CreateCampaign creates a new campaign for a user.
func (s *campaignService) CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error) {

	// Check if user already has a campaign
	if err := s.checkExistingCampaign(userHandle); err != nil {
		return models.Campaign{}, err
	}

	// Validate contributors and get existing/non-existing users
	existing, nonExisting, err := s.authService.FindExistingAndNonExistingUsers(campaign.GetContributorsEmails())
	if err != nil {
		return models.Campaign{}, err
	}

	// Check if existing users can contribute
	var invalidEmails []string
	for _, user := range existing {
		if !user.CanContributeToACampaign() {
			invalidEmails = append(invalidEmails, user.Email)
		}
	}
	if len(invalidEmails) > 0 {
		return models.Campaign{}, errs.BadRequest(
			fmt.Sprintf("Users cannot contribute: %s, already part of another campaign", strings.Join(invalidEmails, ", ")),
			invalidEmails,
		)
	}

	// Create new users for non-existing emails
	if len(nonExisting) > 0 {
		_, err = s.authService.CreateUsers(createUsersFromEmails(nonExisting))
		if err != nil {
			return models.Campaign{}, err
		}
	}

	// Get creator's user details
	user, err := s.authService.GetUserByHandle(userHandle)
	if err != nil {
		return models.Campaign{}, err
	}

	// Setup campaign with creator's details
	campaign.FromBinding(user)

	// Create campaign in database
	campaign, err = s.repo.Create(&campaign)

	if err != nil {
		return models.Campaign{}, errs.InternalServerError(err).Log(s.logger)
	}
	go s.notificationService.NotifyCampaignCreation(&campaign)
	return campaign, nil
}

// TODO: redundant user GetCampaignByIDWithAllRelatedData and select preloads
// GetCampaignByID fetches campaign by ID
func (s *campaignService) GetCampaignByID(id string) (*models.Campaign, error) {
	campaign, err := s.repo.GetByID(id, true)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil, errs.BadRequest("Campaign not found", nil)
		}
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return &campaign, nil
}

// TODO: redundant user GetCampaignByIDWithAllRelatedData and select preloads
// GetCampaignByIDWithContributors fetches campaign by ID with contributors
func (s *campaignService) GetCampaignByIDWithContributors(id string) (*models.Campaign, error) {
	campaign, err := s.repo.GetByIDWithContributors(id)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil, errs.BadRequest("Campaign not found", nil)
		}
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return &campaign, nil
}

// GetCampaignByIDWithAllRelatedData implements interfaces.CampaignService.
func (s *campaignService) GetCampaignByIDWithAllRelatedData(id string) (*models.Campaign, error) {
	preload := models.PreloadOption{
		Images:                 true,
		Payout:                 true,
		Activities:             true,
		ActivitiesContributors: true,
		ActivitiesComments:     true,
		Contributors:           true,
		ContributorsActivities: true,
		CreatedBy:              true,
	}

	campaign, err := s.repo.GetByIDWithSelectedData(id, preload)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil, errs.BadRequest("Campaign not found", nil)
		}
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return &campaign, nil
}

// GetExpiredCampaigns fetches all expired campaigns
func (s *campaignService) GetExpiredCampaigns() ([]models.Campaign, error) {
	//TODO: only admin should be able to get expired campaigns
	campaigns, err := s.repo.GetExpiredCampaigns()
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return campaigns, nil
}

// DeleteCampaign deletes a campaign by ID
func (s *campaignService) DeleteCampaign(campaignID string) error {
	// TODO: only admin should be able to delete campaigns
	if err := s.repo.Delete(campaignID); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

// GetActiveCampaigns fetches all active campaigns
func (s *campaignService) GetActiveCampaigns() ([]models.Campaign, error) {
	campaigns, err := s.repo.GetActiveCampaigns()
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return campaigns, nil
}

// GetNearEndCampaigns fetches all campaigns that are near end
func (s *campaignService) GetNearEndCampaigns() ([]models.Campaign, error) {
	campaigns, err := s.repo.GetNearEndCampaigns()
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return campaigns, nil
}

// Helper Methods

// checkExistingCampaign verifies if user can create a new campaign
func (s *campaignService) checkExistingCampaign(userHandle string) error {
	campaign, err := s.repo.GetByCreatorHandle(userHandle, false)
	if err == nil && campaign.ID != "" {
		return errs.BadRequest("You already have an active campaign", nil)
	}
	if err != nil && !database.Error(err).IsNotfound() {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

// Helper functions

// createUsersFromEmails converts a list of emails to user models
func createUsersFromEmails(emails []string) []models.User {
	users := make([]models.User, 0, len(emails))
	for _, email := range emails {
		users = append(users, *models.NewUser("", email, false))
	}
	return users
}
