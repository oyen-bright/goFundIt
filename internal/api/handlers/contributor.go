package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/contributor"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type ContributorHandler struct {
	service services.ContributorService
}

func NewContributorHandler(service services.ContributorService) *ContributorHandler {
	return &ContributorHandler{
		service: service,
	}
}

// @Summary Add Contributor
// @Description Registers a new contributor to the campaign
// @Tags contributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param request body dto.CreateContributorRequest true "Contributor Details"
// @Success 200 {object} SuccessResponse{data=models.Contributor} "Contributor added to Campaign"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /contributor/{campaignID} [post]
func (h *ContributorHandler) HandleAddContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	campaignKey := getCampaignKey(c)
	var contributor models.Contributor

	//bind request to the contributor model
	if err := c.BindJSON(&contributor); err != nil {
		BadRequest(c, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return
	}

	if err := h.service.AddContributorToCampaign(&contributor, campaignID, campaignKey, claims.Handle); err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Contributor added to Campaign", contributor)

}

// @Summary Remove Contributor
// @Description Removes a contributor from the campaign
// @Tags contributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param contributorID path string true "Contributor ID"
// @Success 200 {object} SuccessResponse "Contributor removed from Campaign"
// @Failure 400 {object} BadRequestResponse "Invalid contributor ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Contributor or Campaign not found"
// @Router /contributor/{campaignID}/{contributorID} [delete]
func (h *ContributorHandler) HandleRemoveContributor(c *gin.Context) {
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)

	contributorId, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid contributor ID", nil)
		return
	}

	if err := h.service.RemoveContributorFromCampaign(contributorId, campaignID, claims.Handle, getCampaignKey(c)); err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Contributor removed from Campaign", nil)

}

// @Summary Edit Contributor
// @Description Modifies a contributor's information
// @Tags contributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param contributorID path string true "Contributor ID"
// @Param request body dto.ContributorEditRequest true "Updated Contributor Details"
// @Success 200 {object} SuccessResponse{data=models.Contributor} "Contributor updated successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Contributor not found"
// @Router /contributor/{campaignID}/{contributorID} [patch]
func (h *ContributorHandler) HandleEditContributor(c *gin.Context) {
	var requestDTO dto.ContributorEditRequest
	var contributor models.Contributor

	claims := getClaimsFromContext(c)

	contributorId, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid contributor ID", nil)
		return
	}
	//bind request to the contributor model
	if err := c.BindJSON(&requestDTO); err != nil {
		BadRequest(c, "Invalid inputs, please check and try again", ExtractValidationErrors(err))
		return
	}
	contributor.Name = requestDTO.Name

	updateContributor, err := h.service.UpdateContributorByID(&contributor, contributorId, claims.Email)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Contributor updated successfully", updateContributor)

}

// @Summary Get All Contributors
// @Description Lists all contributors participating in the campaign
// @Tags contributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} SuccessResponse{data=[]models.Contributor} "Contributors retrieved successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /contributor/{campaignID} [get]
func (h *ContributorHandler) HandleGetContributorsByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	contributors, err := h.service.GetContributorsByCampaignID(campaignID)

	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Contributors retrieved successfully", contributors)

}

// @Summary Get Contributor
// @Description Retrieves detailed information about a specific contributor
// @Tags contributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param contributorID path string true "Contributor ID"
// @Success 200 {object} SuccessResponse{data=models.Contributor} "Contributor retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid contributor ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Contributor not found"
// @Router /contributor/{campaignID}/{contributorID} [get]
func (h *ContributorHandler) HandleGetContributorByID(c *gin.Context) {
	contributorId, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid contributor ID", nil)
		return
	}
	contributor, err := h.service.GetContributorByID(contributorId)

	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Contributor retrieved successfully", contributor)

}
