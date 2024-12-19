package services

import (
	"log"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type contributorService struct {
	repo            repositories.ContributorRepository
	campaignService services.CampaignService
	authService     services.AuthService
	logger          logger.Logger
}

func NewContributorService(repo repositories.ContributorRepository, campaignService services.CampaignService, logger logger.Logger) services.ContributorService {
	return &contributorService{repo: repo, logger: logger, campaignService: campaignService}
}

// AddContributorToCampaign adds a contributor to a campaign
func (s *contributorService) AddContributorToCampaign(contributor *models.Contributor, campaignId, userHandle string) error {

	// Get campaign
	campaign, err := s.campaignService.GetCampaignByID(campaignId)
	if err != nil {
		return err
	}

	// Validate campaign and permissions
	if err := s.validateCampaignAndPermissions(campaign, userHandle, contributor); err != nil {
		return err
	}

	// Create user if it does not exist
	user, err := s.authService.GetUserByEmail(contributor.Email)
	if err != nil {
		return err
	}
	if user.Email == "" {
		newUser := models.NewUser(contributor.Name, contributor.Email, false)
		if err := s.authService.CreateUser(*newUser); err != nil {
			return err
		}
	} else {
		if !user.CanContributeToACampaign() {
			return errs.BadRequest("Not Allowed: contributor is part of another campaigns", nil)
		}
	}

	if err = s.repo.Create(contributor); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

func (s *contributorService) UpdateContributor(contributor *models.Contributor, contributorID uint, userEmail string) (retrievedContributor models.Contributor, err error) {

	// Get contributor
	retrievedContributor, err = s.GetContributorByID(contributorID)
	if err != nil {
		return models.Contributor{}, err
	}

	// Validate ownership
	if retrievedContributor.Email != userEmail {

		log.Println(retrievedContributor.Email, userEmail)
		return models.Contributor{}, errs.BadRequest("Unauthorized: Only contributor can update their details", nil)
	}

	// Update contributor name
	err = s.repo.UpdateName(contributorID, contributor.Name)

	if err != nil {
		return models.Contributor{}, errs.InternalServerError(err).Log(s.logger)
	}

	retrievedContributor.Name = contributor.Name

	return retrievedContributor, nil
}

// RemoveContributorFromCampaign removes a contributor from a campaign
func (s *contributorService) RemoveContributorFromCampaign(contributorId uint, campaignId, userHandle string) error {
	// Get campaign
	campaign, err := s.campaignService.GetCampaignByID(campaignId)
	if err != nil {
		return err
	}

	//validate ownership
	if campaign.CreatedBy.Handle != userHandle {

		log.Println(campaign.CreatedByHandle, userHandle)
		return errs.BadRequest("Unauthorized: Only campaign creator can remove contributors", nil)
	}

	// Get contributor
	contributor := campaign.GetContributorByID(contributorId)
	if contributor.HasPaid() {
		return errs.BadRequest("Cannot remove contributor with paid contribution", nil)
	}

	// Remove contributor
	err = s.repo.Delete(contributor)

	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}
	return nil
}

// GetContributors retrieves contributor by id
func (s *contributorService) GetContributorByID(contributorID uint) (models.Contributor, error) {
	contributor, err := s.repo.GetContributorById(contributorID)

	if err != nil {
		if database.Error(err).IsNotfound() {
			return models.Contributor{}, errs.NotFound("Contributor not found")
		}
		return models.Contributor{}, errs.InternalServerError(err).Log(s.logger)
	}
	return contributor, nil
}

// GetContributors retrieves all contributors in a campaign
func (s *contributorService) GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error) {
	contributors, err := s.repo.GetContributorsByCampaignID(campaignID)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return contributors, nil

}

// Helper methods
func (s *contributorService) validateCampaignAndPermissions(campaign *models.Campaign, userHandle string, contributor *models.Contributor) error {
	// Check creator permission
	if campaign.CreatedBy.Handle != userHandle {
		return errs.BadRequest("Unauthorized: Only campaign creator can add contributors", nil)
	}

	// Check campaign status
	if campaign.HasEnded() {
		return errs.BadRequest("Cannot add contributors: Campaign has ended", nil)
	}

	// Check for existing contributor in this campaign
	if contributor := campaign.GetContributorByEmail(contributor.Email); contributor != nil {
		return errs.BadRequest("Contributor already exists in this campaign", nil)
	}

	return nil
}
