package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
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
		response.BadRequest(c, "Invalid inputs", utils.ExtractValidationErrors(err))
		return
	}

	activity, err := a.service.CreateActivity(activity, claims.Handle, campaignID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity created successfully", activity)
}

// HandleGetActivitiesByCampaignID handles fetching all activities for a campaign
func (a *ActivityHandler) HandleGetActivitiesByCampaignID(c *gin.Context) {
	campaignID := c.Param("campaignID")

	activities, err := a.service.GetActivitiesByCampaignID(campaignID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activities fetched successfully", activities)
}

// HandleGetActivityByID handles fetching a single activity by its ID
func (a *ActivityHandler) HandleGetActivityByID(c *gin.Context) {
	campaignID := c.Param("campaignID")
	activityID, err := parseActivityID(c)
	if err != nil {
		response.BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	activity, err := a.service.GetActivityByID(activityID, campaignID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity fetched successfully", activity)
}

// HandleUpdateActivity handles updating an existing activity
func (a *ActivityHandler) HandleUpdateActivity(c *gin.Context) {
	var activity models.Activity
	claims := getClaimsFromContext(c)
	campaignID := c.Param("campaignID")

	activityID, err := parseActivityID(c)
	if err != nil {
		response.BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := c.BindJSON(&activity); err != nil {
		response.BadRequest(c, "Invalid inputs", utils.ExtractValidationErrors(err))
		return
	}

	activity.ID = activityID
	activity.CampaignID = campaignID

	if err := a.service.UpdateActivity(&activity, claims.Handle); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity updated successfully", activity)
}

// HandleDeleteActivityByID handles deleting an activity
func (a *ActivityHandler) HandleDeleteActivityByID(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := c.Param("campaignID")

	activityID, err := parseActivityID(c)
	if err != nil {
		response.BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := a.service.DeleteActivityByID(activityID, campaignID, claims.Handle); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity deleted successfully", nil)
}

// Helper functions

// getClaimsFromContext extracts JWT claims from the context
func getClaimsFromContext(c *gin.Context) jwt.Claims {
	return c.MustGet("claims").(jwt.Claims)
}

// parseActivityID converts the activity ID from the URL parameter to uint
func parseActivityID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("activityID"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
