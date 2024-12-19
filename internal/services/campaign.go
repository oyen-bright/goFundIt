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
	repo        repositories.CampaignRepository
	authService services.AuthService
	logger      logger.Logger
}

func NewCampaignService(repo repositories.CampaignRepository, authService services.AuthService, logger logger.Logger) services.CampaignService {
	return &campaignService{repo: repo, authService: authService, logger: logger}
}

// CreateCampaign creates a new campaign for a user.
func (s *campaignService) CreateCampaign(campaign models.Campaign, userHandle string) (models.Campaign, error) {
	// Check if user already has a campaign
	if err := s.CheckExistingCampaign(userHandle); err != nil {
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
	return s.repo.Create(&campaign)
}

// createUsersFromEmails converts a list of emails to user models
func createUsersFromEmails(emails []string) []models.User {
	users := make([]models.User, 0, len(emails))
	for _, email := range emails {
		users = append(users, *models.NewUser("", email, false))
	}
	return users
}

// checkExistingCampaign verifies if user can create a new campaign
func (s *campaignService) CheckExistingCampaign(userHandle string) error {
	campaign, err := s.repo.GetByCreatorHandle(userHandle, false)
	if err == nil && campaign.ID != "" {
		return errs.BadRequest("You already have an active campaign", nil)
	}
	if err != nil && !database.Error(err).IsNotfound() {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

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
