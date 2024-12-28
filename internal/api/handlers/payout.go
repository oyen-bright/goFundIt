package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type PayoutHandler struct {
	PayoutService interfaces.PayoutService
}

func NewPayoutHandler(payoutService interfaces.PayoutService) *PayoutHandler {
	return &PayoutHandler{
		PayoutService: payoutService,
	}
}

// GetBankList retrieves a list of banks supported by platform
func (h *PayoutHandler) HandleGetBankList(c *gin.Context) {
	banks, err := h.PayoutService.GetBankList()
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Bank list retrieved successfully", banks)
}

// VerifyAccount verifies an account number and bank code
func (h *PayoutHandler) HandleVerifyAccount(c *gin.Context) {

	var req dto.VerifyAccountRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "Invalid request", ExtractValidationErrors(err))
		return
	}

	account, err := h.PayoutService.VerifyAccount(req)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Account verified successfully", account)
}

// InitializePayout initializes a payout for a campaign
func (h *PayoutHandler) HandleInitializePayout(c *gin.Context) {
	campaignID := GetCampaignID(c)
	claims := getClaimsFromContext(c)
	var req dto.PayoutRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "Invalid request", ExtractValidationErrors(err))
		return
	}

	payout, err := h.PayoutService.InitializePayout(campaignID, claims.Handle, req)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Payout initialized successfully", payout)
}

// HandleInitializeManualPayout initializes a manual payout for a campaign
func (h *PayoutHandler) HandleInitializeManualPayout(c *gin.Context) {
	campaignID := GetCampaignID(c)
	claims := getClaimsFromContext(c)
	var req dto.PayoutRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, "Invalid request", ExtractValidationErrors(err))
		return
	}

	payout, err := h.PayoutService.InitializeManualPayout(campaignID, claims.Handle)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Manual payout initialized successfully", payout)
}

// HandleGetPayoutByCampaignID initialize
func (h *PayoutHandler) HandleGetPayoutByCampaignID(c *gin.Context) {
	campaignID := GetCampaignID(c)

	payout, err := h.PayoutService.GetPayoutByCampaignID(campaignID)

	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Success", payout)
}
