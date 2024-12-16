package activity

import (
	"github.com/oyen-bright/goFundIt/internal/auth"
	"github.com/oyen-bright/goFundIt/internal/campaign"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type ActivityService interface {
	CreateActivity(activity Activity, userHandle, campaignId string) (Activity, error)
	// GetActivitiesByCampaignID(campaignID string) ([]Activity, error)
	// GetActivityByID(activityID uint) (Activity, error)
	// UpdateActivity(activity *Activity) error
	// DeleteActivity(activityID uint) error
	// JoinActivity(contributorID uint, activityID uint) error
	// LeaveActivity(contributorID uint, activityID uint) error
	// GetActivityParticipants(activityID uint) ([]contributor.Contributor, error)
}

type activityService struct {
	repo            ActivityRepository
	authService     auth.AuthService
	campaignService campaign.CampaignService
	logger          logger.Logger
}

func Service(repo ActivityRepository, authService auth.AuthService, campaignService campaign.CampaignService, logger logger.Logger) ActivityService {
	return &activityService{
		repo:            repo,
		authService:     authService,
		campaignService: campaignService,
		logger:          logger,
	}
}

func (a *activityService) CreateActivity(activity Activity, userHandle, campaignId string) (Activity, error) {
	//Check if campaign exist and retrieve campaign data
	campaign, err := a.campaignService.GetCampaignByID(campaignId)
	if err != nil {
		return Activity{}, err

	}

	//check if user exist and part of campaign, approve activity if the user is the campaign creator
	user, err := a.authService.GetUserByHandle(userHandle)
	if err != nil {
		return Activity{}, err
	}

	// if user is the creator of the campaign its approves automatically
	if activity.CreatedByHandle == campaign.CreatedByHandle {
		activity.ApproveActivity()

	} else {
		// only campaign creator can mark an activity as as Mandatory
		activity.MarkAsNotMandatory()
	}
	if isPartOfCampaign := campaign.EmailIsPartOfCampaign(user.Email); !isPartOfCampaign {
		return Activity{}, errs.BadRequest("Sorry, you can't add activities to campaigns you're not part of. Join the campaign to get started!", nil)
	}

	//Update the created by of activity to user data and the campaignId
	activity.UpdateCreatedBy(user)
	activity.UpdateCampaignId(campaign.ID)

	//create the new activity
	activity, err = a.repo.Create(&activity)
	if err != nil {
		return Activity{}, errs.InternalServerError(err).Log(a.logger)
	}

	return activity, nil

}
