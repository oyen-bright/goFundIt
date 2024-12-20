package handlers

import (
	"github.com/gin-gonic/gin"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/response"
)

type SuggestionHandler struct {
	service services.SuggestionService
}

func NewSuggestionHandler(service services.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{service: service}
}

// GetActivitySuggestions generates activity suggestions based on the campaign description.
func (h *SuggestionHandler) HandleGetActivitySuggestions(c *gin.Context) {

	campaignID := GetCampaignID(c)

	suggestions, err := h.service.GetActivitySuggestions(campaignID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity suggestions retrieved successfully", suggestions)
}

// GetActivitySuggestionsViaText generates activity suggestions based on the provided content.
func (h *SuggestionHandler) HandleGetActivitySuggestionsViaText(c *gin.Context) {
	var request struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		response.FromError(c, err)
		return
	}

	suggestions, err := h.service.GetActivitySuggestionsViaText(request.Content)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, "Activity suggestions retrieved successfully", suggestions)
}
