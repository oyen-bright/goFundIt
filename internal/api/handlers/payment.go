package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
)

type PaymentHandler struct {
	service interfaces.PaymentService
}

// NewPaymentHandler creates a new instance of the PaymentHandler
func NewPaymentHandler(service interfaces.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// @Summary Initialize Payment
// @Description Initializes a payment for a contributor
// @Tags payment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param contributorID path string true "Contributor ID"
// @Success 200 {object} SuccessResponse{data=dto.InitializePaymentResponse} "Payment initialized successfully"
// @Failure 400 {object} BadRequestResponse "Invalid contributor ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Contributor not found"
// @Router /payment/contributor/{contributorID} [post]
func (p *PaymentHandler) HandleInitializePayment(c *gin.Context) {
	contributorID, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid contributor ID", nil)
		return
	}

	payment, err := p.service.InitializePayment(contributorID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Payment initialized", payment.GetPaymentLink())
}

// @Summary Initialize Manual Payment
// @Description Initializes a manual payment for a contributor with optional reference file
// @Tags payment
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param contributorID path string true "Contributor ID"
// @Param reference formData file false "Reference file for manual payment"
// @Success 200 {object} SuccessResponse{data=models.Payment} "Manual payment initialized"
// @Failure 400 {object} BadRequestResponse "Invalid request"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Contributor not found"
// @Router /payment/manual/contributor/{contributorID} [post]
func (p *PaymentHandler) HandleInitializeManualPayment(c *gin.Context) {
	userEmail := getClaimsFromContext(c).Email
	contributorID, err := parseContributorID(c)
	if err != nil {
		BadRequest(c, "Invalid contributor ID", nil)
		return
	}
	var reference string

	// Check if reference file was provided
	if file, err := c.FormFile("reference"); err == nil {
		// Only process if a file was actually uploaded
		tmpFile, err := createTempFileFromMultipart(file)
		if err != nil {
			BadRequest(c, "Error processing reference file", err.Error())
			return
		}
		// Clean up the temporary file after we're done
		defer os.Remove(tmpFile.Name())
		reference = tmpFile.Name()
	}

	payment, err := p.service.InitializeManualPayment(contributorID, reference, userEmail, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Manual Payment initialized", payment)
}

// @Summary Verify Payment
// @Description Verifies a payment using the reference
// @Tags payment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param reference path string true "Payment reference"
// @Success 200 {object} SuccessResponse "Payment verified"
// @Failure 400 {object} BadRequestResponse "Invalid reference"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Payment not found"
// @Router /payment/verify/{reference} [post]
func (p *PaymentHandler) HandleVerifyPayment(c *gin.Context) {
	reference := c.Param("reference")
	err := p.service.VerifyPayment(reference)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Payment verified", nil)
}

// @Summary Verify Manual Payment
// @Description Verifies a manual payment using the reference
// @Tags payment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param reference path string true "Payment reference"
// @Success 200 {object} SuccessResponse "Manual payment verified"
// @Failure 400 {object} BadRequestResponse "Invalid reference"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Payment not found"
// @Router /payment/manual/verify/{reference} [post]
func (p *PaymentHandler) HandleVerifyManualPayment(c *gin.Context) {
	reference := c.Param("reference")
	userHandle := getClaimsFromContext(c).Handle

	err := p.service.VerifyManualPayment(reference, userHandle, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Manual Payment verified", nil)
}

// @Summary Handle Paystack Webhook
// @Description Processes incoming Paystack webhook events
// @Tags payment
// @Accept json
// @Produce json
// @Param X-Paystack-Signature header string true "Paystack signature"
// @Param event body paystack.PaystackWebhookEvent true "Webhook event data"
// @Success 200 {object} SuccessResponse "Webhook processed successfully"
// @Failure 400 {object} BadRequestResponse "Invalid webhook data"
// @Router /payment/paystack/webhook [post]
func (p *PaymentHandler) HandlePayStackWebhook(c *gin.Context) {
	// Get the request bo
	var event paystack.PaystackWebhookEvent

	if err := c.BindJSON(&event); err != nil {
		BadRequest(c, "Invalid request", ExtractValidationErrors(err))

	}

	// Handle the event
	p.service.ProcessPaystackWebhook(event)

}
