package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type ActivityHandler struct {
	service services.ActivityService
}

func NewActivityHandler(service services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

// HandleCreateActivity handles the creation of a new activity
func (a *ActivityHandler) HandleCreateActivity(c *gin.Context) {
	var activity models.Activity
	claims := getClaimsFromContext(c)
	campaignID := c.Param("campaignID")

	if err := c.BindJSON(&activity); err != nil {
		BadRequest(c, "Invalid inputs", ExtractValidationErrors(err))
		return
	}

	activity, err := a.service.CreateActivity(activity, claims.Handle, campaignID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity created successfully", activity)
}

// HandleGetActivitiesByCampaignID handles fetching all activities for a campaign
func (a *ActivityHandler) HandleGetActivitiesByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	activities, err := a.service.GetActivitiesByCampaignID(campaignID)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activities fetched successfully", activities)
}

// HandleGetActivityByID handles fetching a single activity by its ID
func (a *ActivityHandler) HandleGetActivityByID(c *gin.Context) {
	campaignID := GetCampaignID(c)
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	activity, err := a.service.GetActivityByID(activityID, campaignID)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity fetched successfully", activity)
}

// HandleApproveActivity handles activity approval
func (a *ActivityHandler) HandleApproveActivity(c *gin.Context) {
	userHandle := getClaimsFromContext(c).Handle

	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	activity, err := a.service.ApproveActivity(activityID, userHandle, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity approved successfully", activity)
}

// HandleUpdateActivity handles updating an existing activity
func (a *ActivityHandler) HandleUpdateActivity(c *gin.Context) {
	var activity models.Activity
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)

	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := c.BindJSON(&activity); err != nil {
		BadRequest(c, "Invalid inputs", ExtractValidationErrors(err))
		return
	}

	activity.ID = activityID
	activity.CampaignID = campaignID

	if err := a.service.UpdateActivity(&activity, claims.Handle); err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity updated successfully", activity)
}

// HandleDeleteActivityByID handles deleting an activity
func (a *ActivityHandler) HandleDeleteActivityByID(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)

	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := a.service.DeleteActivityByID(activityID, campaignID, claims.Handle); err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity deleted successfully", nil)
}

// HandleOptInContributor handles opting in a contributor to an activity
func (a *ActivityHandler) HandleOptInContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	contributorID, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid Contributor ID", nil)
		return
	}
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := a.service.OptInContributor(campaignID, claims.Email, getCampaignKey(c), activityID, contributorID); err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Contributor opted in successfully", nil)
}

// HandleOptOutContributor handles opting out a contributor from an activity
func (a *ActivityHandler) HandleOptOutContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	contributorID, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid Contributor ID", nil)
		return
	}
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := a.service.OptOutContributor(campaignID, claims.Email, getCampaignKey(c), activityID, contributorID); err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Contributor opted out successfully", nil)
}

// HandleGetParticipants handles fetching all participants for an activity
func (a *ActivityHandler) HandleGetParticipants(c *gin.Context) {
	campaignID := GetCampaignID(c)
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	participants, err := a.service.GetParticipants(activityID, campaignID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Participants fetched successfully", participants)
}
