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

// InitializePayment initializes a payment for a contributor
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

// InitializeManualPayment initializes a manual payment for a contributor
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

// VerifyPayment verifies a payment
func (p *PaymentHandler) HandleVerifyPayment(c *gin.Context) {
	reference := c.Param("reference")
	err := p.service.VerifyPayment(reference)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Payment verified", nil)
}

// VerifyManualPayment verifies a manual payment
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

// HandleWebhook handles incoming payment webhooks
func (p *PaymentHandler) HandlePayStackWebhook(c *gin.Context) {
	// Get the request bo
	var event paystack.PaystackWebhookEvent

	if err := c.BindJSON(&event); err != nil {
		BadRequest(c, "Invalid request", ExtractValidationErrors(err))

	}

	// Handle the event
	p.service.ProcessPaystackWebhook(event)

}
