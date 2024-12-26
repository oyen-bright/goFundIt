package handlers

import (
	"github.com/gin-gonic/gin"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/response"
)

type AnalyticsHandle struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandle {
	return &AnalyticsHandle{
		service: service,
	}
}

func (h *AnalyticsHandle) HandleProcessAnalyticsNow(c *gin.Context) {
	err := h.service.ProcessAnalyticsNow()
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Analytics proceed and send to email", nil)
}
