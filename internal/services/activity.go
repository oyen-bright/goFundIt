package services

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type activityService struct {
	repo            repositories.ActivityRepository
	authService     services.AuthService
	campaignService services.CampaignService
	logger          logger.Logger
}

func NewActivityService(repo repositories.ActivityRepository, authService services.AuthService, campaignService services.CampaignService, logger logger.Logger) services.ActivityService {
	return &activityService{
		repo:            repo,
		authService:     authService,
		campaignService: campaignService,
		logger:          logger,
	}
}

func (a *activityService) CreateActivity(activity models.Activity, userHandle, campaignId string) (models.Activity, error) {
	//Check if campaign exist and retrieve campaign data
	campaign, err := a.campaignService.GetCampaignByID(campaignId)
	if err != nil {
		return models.Activity{}, err

	}

	//check if user exist and part of campaign, approve activity if the user is the campaign creator
	user, err := a.authService.GetUserByHandle(userHandle)
	if err != nil {
		return models.Activity{}, err
	}

	// if user is the creator of the campaign its approves automatically
	if activity.CreatedByHandle == campaign.CreatedByHandle {
		activity.ApproveActivity()

	} else {
		// only campaign creator can mark an activity as as Mandatory
		activity.MarkAsNotMandatory()
	}

	if isPartOfCampaign := campaign.EmailIsPartOfCampaign(user.Email); !isPartOfCampaign {
		return models.Activity{}, errs.BadRequest("Sorry, you can't add activities to campaigns you're not part of. Join the campaign to get started!", nil)
	}

	//Update the created by of activity to user data and the campaignId
	activity.UpdateCreatedBy(user)
	activity.UpdateCampaignId(campaign.ID)

	//create the new activity
	activity, err = a.repo.Create(&activity)
	if err != nil {
		return models.Activity{}, errs.InternalServerError(err).Log(a.logger)
	}

	return activity, nil

}
