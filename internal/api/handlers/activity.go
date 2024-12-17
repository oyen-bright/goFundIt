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

// handleNewActivity handles incoming activity creation requests.
// Required Fields:
//   - Title: required, string
//   - Subtitle: optional, string
//   - imageUrl: optional, string
//   - isMandatory; optional, boolean, defaults:false
func (a *ActivityHandler) HandleCreateActivity(context *gin.Context) {

	var activity models.Activity

	claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")

	//bind request to the Activity model
	if err := context.BindJSON(&activity); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	//create activity
	activity, err := a.service.CreateActivity(activity, claims.Handle, campaignID)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activity created successfully", activity)

}

func (a *ActivityHandler) HandleGetActivitiesByCampaignID(context *gin.Context) {

	// claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")

	//get Activities
	activities, err := a.service.GetActivitiesByCampaignID(campaignID)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activities Fetched successfully", activities)

}

func (a *ActivityHandler) HandleGetActivityByID(context *gin.Context) {

	// claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")
	activityID, err := strconv.ParseUint(context.Param("activityID"), 10, 64)

	if err != nil {
		response.BadRequest(context, "Invalid Activity ID", nil)
		return
	}

	//get Activity
	activities, err := a.service.GetActivityByID(uint(activityID), campaignID)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activity Fetched successfully", activities)

}

func (a *ActivityHandler) HandleUpdateActivity(context *gin.Context) {
	var activity models.Activity

	claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")
	activityID, err := strconv.ParseUint(context.Param("activityID"), 10, 64)

	if err != nil {
		response.BadRequest(context, "Invalid Activity ID", nil)
		return
	}

	//bind request to the Activity model
	if err := context.BindJSON(&activity); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	activity.ID = uint(activityID)
	activity.CampaignID = campaignID

	//get Activity
	err = a.service.UpdateActivity(&activity, claims.Handle)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activity Updated successfully", activity)

}

func (a *ActivityHandler) HandleDeleteActivityByID(context *gin.Context) {

	claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")
	activityID, err := strconv.ParseUint(context.Param("activityID"), 10, 64)

	if err != nil {
		response.BadRequest(context, "Invalid Activity ID", nil)
		return
	}

	//get Activity
	err = a.service.DeleteActivityByID(uint(activityID), campaignID, claims.Handle)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activity deleted successfully", nil)

}
