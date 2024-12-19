package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
)

type CampaignHandler struct {
	service services.CampaignService
}

func NewCampaignHandler(service services.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		service: service,
	}
}

// handleCreateCampaign handles incoming campaign creation requests.
// Required Fields:
//   - Title: required, string
//   - Description: optional, string
//   - Images: []{imageUrl: required, url}
//   - Activities: []{title: required, string, subtitle: optional, string, imageUrl: optional, string, isMandatory: optional, bool}
//   - TargetAmount: required, number
//   - StartDate: required, date
//   - EndDate: required, date
func (h *CampaignHandler) HandleCreateCampaign(c *gin.Context) {

	claims := c.MustGet("claims").(jwt.Claims)
	var campaign models.Campaign

	//bind request to the campaign model
	if err := c.BindJSON(&campaign); err != nil {
		response.BadRequest(c, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	//validate the request
	if err := campaign.ValidateNewCampaign(); err != nil {
		response.BadRequest(c, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	//create service
	campaign, err := h.service.CreateCampaign(campaign, claims.Handle)

	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Campaigns created successfully", campaign)

}

func (h *CampaignHandler) HandleGetCampaignByID(c *gin.Context) {

	campaignID := c.Param("id")

	campaign, err := h.service.GetCampaignByID(campaignID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Campaigns retrieved successfully", campaign)

}

func (c *CampaignHandler) HandleUpdateCampaignByID(context *gin.Context) {

	response.Success(context, "Campaigns Updated successfully", []string{})

}
