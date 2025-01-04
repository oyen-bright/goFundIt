package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/campaign"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

type CampaignHandler struct {
	service services.CampaignService
}

func NewCampaignHandler(service services.CampaignService) *CampaignHandler {
	return &CampaignHandler{service: service}
}

// @Summary Create Campaign
// @Description Creates a new campaign with the provided details
// @Tags campaign
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param request body dto.CampaignRequest true "Campaign Details"
// @Success 200 {object} SuccessResponse{data=models.Campaign} "Campaign created successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs, please check and try again"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /campaign/create [post]
func (h *CampaignHandler) HandleCreateCampaign(c *gin.Context) {
	userHandle := getUserHandle(c)
	//TODO: use dto.CampaignRequest
	var campaign models.Campaign
	if err := bindJSON(c, &campaign); err != nil {
		return
	}

	createdCampaign, err := h.service.CreateCampaign(&campaign, userHandle)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Campaign created successfully", createdCampaign)
}

// @Summary Get Campaign
// @Description Retrieves a campaign by its ID
// @Tags campaign
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} SuccessResponse{data=models.Campaign} "Campaign retrieved successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid campaign ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /campaign/{campaignID} [get]
func (h *CampaignHandler) HandleGetCampaignByID(c *gin.Context) {
	campaignID := GetCampaignID(c)
	campaign, err := h.service.GetCampaignByID(campaignID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Campaign retrieved successfully", campaign)
}

// @Summary Update Campaign
// @Description Updates an existing campaign by its ID
// @Tags campaign
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param request body dto.CampaignUpdateRequest true "Update Campaign Details"
// @Success 200 {object} SuccessResponse{data=models.Campaign} "Campaign updated successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /campaign/{campaignID} [patch]
func (h *CampaignHandler) HandleUpdateCampaignByID(c *gin.Context) {
	var requestDTO dto.CampaignUpdateRequest
	userHandle := getUserHandle(c)
	campaignID := GetCampaignID(c)

	if err := bindJSON(c, &requestDTO); err != nil {
		return
	}

	campaign, err := h.service.UpdateCampaign(requestDTO, campaignID, userHandle, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Campaign updated successfully", campaign)
}

// Helper Functions -----------------------------------------------------------------

func getUserHandle(c *gin.Context) string {
	return c.MustGet("claims").(jwt.Claims).Handle
}
