package services

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

type contributorService struct {
	repo                repositories.ContributorRepository
	campaignService     services.CampaignService
	analyticsService    services.AnalyticsService
	authService         services.AuthService
	notificationService services.NotificationService
	broadcaster         services.EventBroadcaster
	logger              logger.Logger
	runAsync            func(func())
}

func NewContributorService(
	repo repositories.ContributorRepository,
	campaignService services.CampaignService,
	analyticsService services.AnalyticsService,
	authService services.AuthService,
	notificationService services.NotificationService,
	broadcaster services.EventBroadcaster,
	logger logger.Logger,
) services.ContributorService {
	return &contributorService{
		repo:                repo,
		campaignService:     campaignService,
		analyticsService:    analyticsService,
		authService:         authService,
		notificationService: notificationService,
		broadcaster:         broadcaster,
		logger:              logger,
		runAsync:            func(f func()) { go f() },
	}
}

// AddContributorToCampaign adds a contributor to a campaign
func (s *contributorService) AddContributorToCampaign(contributor *models.Contributor, campaignId, campaignKey, userHandle string) error {

	// Get campaign
	campaign, err := s.campaignService.GetCampaignByID(campaignId, campaignKey)
	if err != nil {
		return err
	}

	// Validate campaign and permissions
	if err := s.validateCampaignAndPermissions(campaign, userHandle, contributor); err != nil {
		return err
	}

	// Create user if it does not exist
	user, err := s.authService.FindUserByEmail(contributor.Email)
	if err != nil {
		return err
	}
	if user == nil {
		newUser := models.NewUser(contributor.Name, contributor.Email, false)
		if err := s.authService.CreateUser(*newUser); err != nil {
			return err
		}
	} else {
		if !user.CanContributeToACampaign() {
			return errs.BadRequest("Not Allowed: contributor is part of another campaigns", nil)
		}
	}

	contributor.CampaignID = campaignId
	campaign.Key = campaignKey

	if err = s.repo.Create(contributor); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	// broadcast event
	// go s.broadcaster.NewEvent(campaign.ID, websocket.EventTypeContributionCreated, contributor)
	// // send notification
	// go s.notificationService.NotifyContributorAdded(contributor, campaign)
	// // calculates the new target amount and broadcast event
	// go s.campaignService.RecalculateTargetAmount(campaignId)

	//
	s.runAsync(func() {
		s.broadcaster.NewEvent(campaign.ID, websocket.EventTypeContributionCreated, contributor)
		s.notificationService.NotifyContributorAdded(contributor, campaign)
		s.campaignService.RecalculateTargetAmount(campaignId)
	})

	return nil
}

func (s *contributorService) UpdateContributor(contributor *models.Contributor) error {

	// Update contributor
	err := s.repo.Update(contributor)

	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	// broadcast event
	// go s.broadcaster.NewEvent(contributor.CampaignID, websocket.EventTypeContributorUpdated, contributor)
	// calculates the new target amount and broadcast event
	// go s.campaignService.RecalculateTargetAmount(contributor.CampaignID)
	s.runAsync(func() {
		s.broadcaster.NewEvent(contributor.CampaignID, websocket.EventTypeContributorUpdated, contributor)
		s.campaignService.RecalculateTargetAmount(contributor.CampaignID)
	})

	return nil
}

func (s *contributorService) UpdateContributorByID(contributor *models.Contributor, contributorID uint, userEmail string) (retrievedContributor models.Contributor, err error) {

	// Get contributor
	retrievedContributor, err = s.GetContributorByID(contributorID)
	if err != nil {
		return models.Contributor{}, err
	}

	// Validate ownership
	if retrievedContributor.Email != userEmail {

		return models.Contributor{}, errs.BadRequest("Unauthorized: Only contributor can update their details", nil)
	}

	// Update contributor name
	err = s.repo.UpdateName(contributorID, contributor.Name)

	if err != nil {
		return models.Contributor{}, errs.InternalServerError(err).Log(s.logger)
	}

	retrievedContributor.Name = contributor.Name

	// broadcast event
	go s.broadcaster.NewEvent(retrievedContributor.CampaignID, websocket.EventTypeContributorUpdated, retrievedContributor)

	return retrievedContributor, nil
}

// RemoveContributorFromCampaign removes a contributor from a campaign
func (s *contributorService) RemoveContributorFromCampaign(contributorId uint, campaignId, userHandle, key string) error {
	// Get campaign
	campaign, err := s.campaignService.GetCampaignByID(campaignId, key)
	if err != nil {
		return err
	}

	//validate ownership
	if campaign.CreatedBy.Handle != userHandle {

		return errs.BadRequest("Unauthorized: Only campaign creator can remove contributors", nil)
	}

	// Get contributor
	contributor := campaign.GetContributorByID(contributorId)
	if contributor == nil {
		return errs.NotFound("Contributor not found")
	}
	if contributor.HasPaid() {
		return errs.BadRequest("Cannot remove contributor with paid contribution", nil)
	}

	// Remove contributor
	err = s.repo.Delete(contributor)

	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	// broadcast event
	// go s.broadcaster.NewEvent(campaign.ID, websocket.EventTypeContributorDeleted, contributor)

	// calculates the new target amount and broadcast event
	// go s.campaignService.RecalculateTargetAmount(contributor.CampaignID)

	s.runAsync(func() {
		s.broadcaster.NewEvent(campaign.ID, websocket.EventTypeContributorDeleted, contributor)
		s.campaignService.RecalculateTargetAmount(contributor.CampaignID)

	})

	return nil
}

// GetContributors retrieves contributor by id
func (s *contributorService) GetContributorByID(contributorID uint) (models.Contributor, error) {
	contributor, err := s.repo.GetContributorById(contributorID, true)

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

// Helper Methods --------------------------------------

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
