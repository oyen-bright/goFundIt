package handlers

import (
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
func (a *ActivityHandler) HandleNewActivity(context *gin.Context) {

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
	response.Success(context, "Activity created successfully", activity.ToJSON())

}
