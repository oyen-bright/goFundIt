package activity

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/interfaces"
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	"github.com/oyen-bright/goFundIt/internal/utils/response"
)

type activityHandler struct {
	service ActivityService
}

func Handler(service ActivityService) interfaces.HandlerInterface {
	return &activityHandler{
		service: service,
	}
}

func (a *activityHandler) RegisterRoutes(activityRoutes *gin.RouterGroup, middlewares []gin.HandlerFunc) {

	activityRoutes.Use(middlewares...)
	activityRoutes.POST("/:campaignID/new", a.handleNewActivity)
}

// handleNewActivity handles incoming activity creation requests.
// Required Fields:
//   - Title: required, string
//   - Subtitle: optional, string
//   - imageUrl: optional, string
//   - isMandatory; optional, boolean, defaults:false
func (a *activityHandler) handleNewActivity(context *gin.Context) {

	var activity Activity

	claims := context.MustGet("claims").(jwt.Claims)
	campaignID := context.Param("campaignID")

	//bind request to the Activity model
	if err := context.BindJSON(&activity); err != nil {
		response.BadRequest(context, "Invalid inputs, please check and try again", response.ExtractValidationErrors(err))
		return
	}

	//create activity
	activity, err := a.service.CreateActivity(activity, claims.Handle, campaignID)

	if err != nil {
		response.FromError(context, err)
		return
	}
	response.Success(context, "Activity created successfully", activity.ToJSON())

}
