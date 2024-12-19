package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/contributor"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
)

type ContributorHandler struct {
	service services.ContributorService
}

func NewContributorHandler(service services.ContributorService) *ContributorHandler {
	return &ContributorHandler{
		service: service,
	}
}

// HandleAddContributor handles adding contributor to a campaign
func (h *ContributorHandler) HandleAddContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	var contributor models.Contributor

	//bind request to the contributor model
	if err := c.BindJSON(&contributor); err != nil {
		response.BadRequest(c, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}

	if err := h.service.AddContributorToCampaign(&contributor, campaignID, claims.Handle); err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Contributor added to Campaign", contributor)

}

// HandleRemoveContributor handles removing a contributor from a campaign
func (h *ContributorHandler) HandleRemoveContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)

	contributorId, err := parseContributorID(c)
	if err != nil {
		response.BadRequest(c, "Invalid contributor ID", nil)
		return
	}

	if err := h.service.RemoveContributorFromCampaign(contributorId, campaignID, claims.Handle); err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Contributor removed from Campaign", nil)

}

// HandleEditContributor handles editing contributor data
func (h *ContributorHandler) HandleEditContributor(c *gin.Context) {
	var requestDTO dto.ContributorEditRequest
	var contributor models.Contributor

	claims := getClaimsFromContext(c)

	contributorId, err := parseContributorID(c)
	if err != nil {
		response.BadRequest(c, "Invalid contributor ID", nil)
		return
	}
	//bind request to the contributor model
	if err := c.BindJSON(&requestDTO); err != nil {
		response.BadRequest(c, "Invalid inputs, please check and try again", utils.ExtractValidationErrors(err))
		return
	}
	contributor.Name = requestDTO.Name

	updateContributor, err := h.service.UpdateContributor(&contributor, contributorId, claims.Email)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Contributor updated successfully", updateContributor)

}

// HandleGetContributorsByCampaignID handles fetching all contributors to a campaign
func (h *ContributorHandler) HandleGetContributorsByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	contributors, err := h.service.GetContributorsByCampaignID(campaignID)

	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Contributors retrieved successfully", contributors)

}

// HandleGetContributorByID handles fetching a contributor by the contributor ID
func (h *ContributorHandler) HandleGetContributorByID(c *gin.Context) {
	contributorId, err := parseContributorID(c)
	if err != nil {
		response.BadRequest(c, "Invalid contributor ID", nil)
		return
	}
	contributor, err := h.service.GetContributorByID(contributorId)

	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Contributor retrieved successfully", contributor)

}
