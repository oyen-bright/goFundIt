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

// @Summary Create Activity
// @Description Creates a new activity for a campaign
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param request body dto.ActivityRequest true "Activity Details"
// @Success 200 {object} SuccessResponse{data=models.Activity} "Activity created successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID} [post]
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

// @Summary Get Activities
// @Description Retrieves all activities for a campaign
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} SuccessResponse{data=[]models.Activity} "Activities fetched successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /activity/{campaignID} [get]
func (a *ActivityHandler) HandleGetActivitiesByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	activities, err := a.service.GetActivitiesByCampaignID(campaignID)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activities fetched successfully", activities)
}

// @Summary Get Activity
// @Description Retrieves a specific activity by ID
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Success 200 {object} SuccessResponse{data=models.Activity} "Activity fetched successfully"
// @Failure 400 {object} BadRequestResponse "Invalid Activity ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Activity not found"
// @Router /activity/{campaignID}/{activityID} [get]
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

// @Summary Approve Activity
// @Description Approves an activity in a campaign
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Success 200 {object} SuccessResponse{data=models.Activity} "Activity approved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid Activity ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/approve [post]
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

// @Summary Update Activity
// @Description Updates an existing activity
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param request body dto.UpdateActivityRequest true "Update Activity Details"
// @Success 200 {object} SuccessResponse{data=models.Activity} "Activity updated successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID} [patch]
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

// @Summary Delete Activity
// @Description Deletes an activity from a campaign
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Success 200 {object} SuccessResponse "Activity deleted successfully"
// @Failure 400 {object} BadRequestResponse "Invalid Activity ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID} [delete]
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

// @Summary Opt In Contributor
// @Description Opts in a contributor to an activity
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param contributorID path string true "Contributor ID"
// @Success 200 {object} SuccessResponse "Contributor opted in successfully"
// @Failure 400 {object} BadRequestResponse "Invalid IDs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/participants/{contributorID} [post]
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

// @Summary Opt Out Contributor
// @Description Opts out a contributor from an activity
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param contributorID path string true "Contributor ID"
// @Success 200 {object} SuccessResponse "Contributor opted out successfully"
// @Failure 400 {object} BadRequestResponse "Invalid IDs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/participants/{contributorID} [delete]
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

// @Summary Get Participants
// @Description Retrieves all participants for an activity
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Success 200 {object} SuccessResponse{data=[]models.Contributor} "Participants fetched successfully"
// @Failure 400 {object} BadRequestResponse "Invalid Activity ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/participants [get]
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
