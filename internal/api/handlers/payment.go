package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
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
		response.BadRequest(c, "Invalid contributor ID", nil)
		return
	}

	payment, err := p.service.InitializePayment(contributorID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Payment initialized", payment.GetPaymentLink())
}

// VerifyPayment verifies a payment
func (p *PaymentHandler) HandleVerifyPayment(c *gin.Context) {
	reference := c.Param("reference")
	err := p.service.VerifyPayment(reference)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Payment verified", nil)
}

// HandleWebhook handles incoming payment webhooks
func (p *PaymentHandler) HandlePayStackWebhook(c *gin.Context) {
	// Get the request bo
	var event paystack.PaystackWebhookEvent

	if err := c.BindJSON(&event); err != nil {
		response.BadRequest(c, "Invalid request", utils.ExtractValidationErrors(err))
		return
	}

	// Handle the event
	p.service.HandlePaystackWebhook(event)

}
