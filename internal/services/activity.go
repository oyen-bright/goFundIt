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

type activityService struct {
	repo            repositories.ActivityRepository
	authService     services.AuthService
	campaignService services.CampaignService
	logger          logger.Logger
}

func NewActivityService(
	repo repositories.ActivityRepository,
	authService services.AuthService,
	campaignService services.CampaignService,
	logger logger.Logger,
) services.ActivityService {
	return &activityService{
		repo:            repo,
		authService:     authService,
		campaignService: campaignService,
		logger:          logger,
	}
}

// CreateActivity handles the creation of a new activity
func (s *activityService) CreateActivity(activity models.Activity, userHandle, campaignID string) (models.Activity, error) {
	// Validate campaign and user
	campaign, user, err := s.validateCampaignAndUser(campaignID, userHandle)
	if err != nil {
		return models.Activity{}, err
	}

	// Check if user is part of campaign
	if !campaign.EmailIsPartOfCampaign(user.Email) {
		return models.Activity{}, errs.BadRequest(
			"Sorry, you can't add activities to campaigns you're not part of. Join the campaign to get started!",
			nil,
		)
	}

	// Setup activity
	s.setupActivity(&activity, campaign, user)

	// Create activity
	createdActivity, err := s.repo.Create(&activity)
	if err != nil {
		return models.Activity{}, (errs.InternalServerError(err)).Log(s.logger)
	}

	return createdActivity, nil
}

// GetActivitiesByCampaignID retrieves all activities for a campaign
func (s *activityService) GetActivitiesByCampaignID(campaignID string) ([]models.Activity, error) {
	activities, err := s.repo.GetActivitiesByCampaignID(campaignID)
	if err != nil {
		return nil, (errs.InternalServerError(err)).Log(s.logger)
	}
	return activities, nil
}

// GetActivityByID retrieves a specific activity
func (s *activityService) GetActivityByID(activityID uint, campaignID string) (models.Activity, error) {
	activity, err := s.repo.GetActivityByID(activityID)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return models.Activity{}, (errs.NotFound("Activity not found")).Log(s.logger)
		}
		return models.Activity{}, (errs.InternalServerError(err)).Log(s.logger)
	}

	if activity.CampaignID != campaignID {
		return models.Activity{}, (errs.NotFound("Activity does not belong to this campaign")).Log(s.logger)
	}

	return activity, nil
}

// UpdateActivity updates an existing activity
func (s *activityService) UpdateActivity(activity *models.Activity, userHandle string) error {
	existingActivity, err := s.validateActivityForModification(
		activity.ID,
		activity.CampaignID,
		userHandle,
	)
	if err != nil {
		return err
	}

	if existingActivity.GetPaidContributorsCount() > 0 {
		return errs.BadRequest("Cannot update activity with paid contributors", activity)
	}

	if err := s.repo.Update(activity); err != nil {
		return (errs.InternalServerError(err)).Log(s.logger)
	}

	return nil
}

// DeleteActivityByID deletes an activity
func (s *activityService) DeleteActivityByID(activityID uint, campaignID, userHandle string) error {
	activity, err := s.validateActivityForModification(activityID, campaignID, userHandle)
	if err != nil {
		return err
	}

	if activity.GetPaidContributorsCount() > 0 {
		return errs.BadRequest("Cannot delete activity with paid contributors", activity)
	}

	if err := s.repo.Delete(&activity); err != nil {
		return (errs.InternalServerError(err)).Log(s.logger)
	}

	return nil
}

// OptInContributor opts in a contributor to an activity
func (s *activityService) OptInContributor(campaignID, userEmail string, activityID, contributorID uint) error {

	// Validate campaign and user
	campaign, err := s.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		return err
	}

	//Validate Contributor and Activity
	contributor, activity, err := s.validateContributorActivityForOptInOptOut(campaign, contributorID, activityID, userEmail)
	if err != nil {
		return err
	}

	if activity.IsContributorOptedIn(contributorID) {
		return errs.BadRequest("Contributor has already opted in.", nil)
	}

	if err := s.repo.AddContributorToActivity(activity.ID, contributor.ID); err != nil {
		return (errs.InternalServerError(err)).Log(s.logger)
	}

	return nil
}

// OptOutContributor opts out a contributor from an activity
func (s *activityService) OptOutContributor(campaignID, userEmail string, activityID, contributorID uint) error {
	// Validate campaign and user
	campaign, err := s.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		return err
	}

	//Validate Contributor and Activity
	contributor, activity, err := s.validateContributorActivityForOptInOptOut(campaign, contributorID, activityID, userEmail)
	if err != nil {
		return err
	}
	if !activity.IsContributorOptedIn(contributorID) {
		return errs.BadRequest("Contributor has already opted out.", nil)
	}

	if err := s.repo.RemoveContributorFromActivity(activityID, contributor.ID); err != nil {
		return (errs.InternalServerError(err)).Log(s.logger)
	}

	return nil
}

// GetParticipants retrieves all contributors for a specific activity
func (s *activityService) GetParticipants(activityID uint, campaignID string) ([]models.Contributor, error) {

	// Validate campaign and user
	campaign, err := s.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		return nil, err
	}
	if activity := campaign.GetActivityById(activityID); activity == nil {
		return nil, errs.BadRequest("Activity not found in this campaign.", activityID)
	}

	contributors, err := s.repo.GetActivityParticipants(activityID)

	if err != nil {
		return nil, (errs.InternalServerError(err)).Log(s.logger)
	}
	return contributors, nil
}

// Helper methods

func (s *activityService) validateCampaignAndUser(campaignID, userHandle string) (*models.Campaign, *models.User, error) {
	campaign, err := s.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.authService.GetUserByHandle(userHandle)
	if err != nil {
		return nil, nil, err
	}

	return campaign, &user, nil
}

func (s *activityService) setupActivity(activity *models.Activity, campaign *models.Campaign, user *models.User) {
	activity.UpdateCreatedBy(*user)
	activity.UpdateCampaignId(campaign.ID)

	if activity.CreatedByHandle == campaign.CreatedByHandle {
		activity.ApproveActivity()
	} else {
		activity.MarkAsNotMandatory()
	}
}

func (s *activityService) validateActivityForModification(activityID uint, campaignID, userHandle string) (models.Activity, error) {
	activity, err := s.repo.GetActivityByID(activityID)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return models.Activity{}, (errs.NotFound("Activity not found")).Log(s.logger)
		}
		return models.Activity{}, (errs.InternalServerError(err)).Log(s.logger)
	}

	if activity.CampaignID != campaignID {
		return models.Activity{}, errs.BadRequest("Activity does not belong to Campaign", activity)
	}

	if userHandle != activity.CreatedByHandle {
		return models.Activity{}, errs.Forbidden("You are not authorized to modify this activity")
	}

	return activity, nil
}

func (s *activityService) validateContributorActivityForOptInOptOut(campaign *models.Campaign, contributorID, activityID uint, userEmail string) (contributor *models.Contributor, activity *models.Activity, err error) {
	// Validate Contributor
	contributor = campaign.GetContributorByID(contributorID)
	if contributor == nil {
		return nil, nil, errs.BadRequest("Contributor not found in this campaign.", contributorID)
	}
	if contributor.Email != userEmail {
		return nil, nil, errs.BadRequest("Only the contributor can perform this action.", nil)
	}
	if contributor.HasPaid() {
		return nil, nil, errs.BadRequest("Action cannot be performed after making a payment.", nil)
	}
	log.Println(len(campaign.Activities))
	// Validate activity
	activity = campaign.GetActivityById(activityID)
	if activity == nil {
		return nil, nil, errs.BadRequest("Activity not found in this campaign.", activityID)
	}
	if !activity.IsApproved {
		return nil, nil, errs.BadRequest("Activity is not approved.", activity)
	}

	return contributor, activity, nil
}
