package handlers

import (
	"github.com/gin-gonic/gin"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type AnalyticsHandler struct {
	service services.AnalyticsService
}

// NewAnalyticsHandler creates a new instance of AnalyticsHandle
func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// @Summary Process Analytics Now
// @Description Triggers immediate processing of analytics data and sends results via email
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Success 200 {object} SuccessResponse "Analytics processed and sent to email"
// @Failure 400 {object} BadRequestResponse "Processing failed"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 500 {object} response "Internal server error"
// @Router /analytics/process [get]
func (h *AnalyticsHandler) HandleProcessAnalyticsNow(c *gin.Context) {
	err := h.service.ProcessAnalyticsNow()
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Analytics processed and sent to email", nil)
}
