package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type PayoutHandler struct {
	service services.PayoutService
}

func NewPayoutHandler(service services.PayoutService) *PayoutHandler {
	return &PayoutHandler{
		service: service,
	}
}

// @Summary Get Bank List
// @Description Retrieves list of available banks for payout
// @Tags payout
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse{data=[]paystack.Bank} "Bank list retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid request"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /payout/bank-list [get]
func (p *PayoutHandler) HandleGetBankList(c *gin.Context) {
	banks, err := p.service.GetBankList()
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Bank list retrieved successfully", banks)
}

// @Summary Verify Bank Account
// @Description Verifies a bank account for payout
// @Tags payout
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.VerifyAccountRequest true "Account verification details"
// @Success 200 {object} SuccessResponse{data=paystack.ResolveAccountResponse} "Account verified successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid account details"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /payout/verify/bank-account [post]
func (p *PayoutHandler) HandleVerifyAccount(c *gin.Context) {
	var req dto.VerifyAccountRequest
	if err := bindJSON(c, &req); err != nil {
		return
	}

	account, err := p.service.VerifyAccount(req)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Account verified successfully", account)
}

// @Summary Initialize Payout
// @Description Initializes a payout for a campaign
// @Tags payout
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param request body dto.PayoutRequest true "Payout details"
// @Success 200 {object} SuccessResponse{data=models.Payout} "Payout initialized successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid payout details"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /payout/{campaignID} [post]
func (p *PayoutHandler) HandleInitializePayout(c *gin.Context) {
	var req dto.PayoutRequest
	if err := bindJSON(c, &req); err != nil {
		return
	}

	campaignID := GetCampaignID(c)
	userHandle := getClaimsFromContext(c).Handle

	payout, err := p.service.InitializePayout(campaignID, userHandle, req)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Payout initialized successfully", payout)
}

// @Summary Initialize Manual Payout
// @Description Initializes a manual payout for a campaign
// @Tags payout
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param request body dto.PayoutRequest true "Manual payout details"
// @Success 200 {object} SuccessResponse{data=models.Payout} "Manual payout initialized successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid payout details"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /payout/manual/{campaignID} [post]
func (p *PayoutHandler) HandleInitializeManualPayout(c *gin.Context) {

	campaignID := GetCampaignID(c)
	userHandle := getClaimsFromContext(c).Handle

	payout, err := p.service.InitializeManualPayout(campaignID, userHandle)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Manual payout initialized successfully", payout)
}

// @Summary Get Campaign Payout
// @Description Retrieves payout information for a campaign
// @Tags payout
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} SuccessResponse{data=models.Payout} "Payout information retrieved successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /payout/{campaignID} [get]
func (p *PayoutHandler) HandleGetPayoutByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	payout, err := p.service.GetPayoutByCampaignID(campaignID)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Payout information retrieved successfully", payout)
}
