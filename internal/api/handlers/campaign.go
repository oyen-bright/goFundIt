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

// HandleCreateCampaign handles the creation of a new campaign
func (h *CampaignHandler) HandleCreateCampaign(c *gin.Context) {
	userHandle := getUserHandle(c)
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

// HandleGetCampaignByID handles fetching a campaign by ID
func (h *CampaignHandler) HandleGetCampaignByID(c *gin.Context) {
	campaignID := GetCampaignID(c)
	campaign, err := h.service.GetCampaignByID(campaignID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Campaign retrieved successfully", campaign)
}

// HandleUpdateCampaignByID handles updating a campaign by ID
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
