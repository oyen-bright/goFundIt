package handlers

import (
	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/suggestion"

	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type SuggestionHandler struct {
	service services.SuggestionService
}

// NewSuggestionHandler creates a new instance of SuggestionHandler
func NewSuggestionHandler(service services.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{service: service}
}

// @Summary Get Campaign Activity Suggestions
// @Description Retrieves a list of suggested activities for a specific campaign
// @Tags suggestion
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} SuccessResponse{data=[]models.ActivitySuggestion} "Activity suggestions retrieved successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /suggestions/activity/{campaignID} [get]
func (h *SuggestionHandler) HandleGetActivitySuggestions(c *gin.Context) {
	campaignID := GetCampaignID(c)

	suggestions, err := h.service.GetActivitySuggestions(campaignID, getCampaignKey(c))
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity suggestions retrieved successfully", suggestions)
}

// @Summary Get Activity Suggestions by Text
// @Description Generates activity suggestions based on provided content text
// @Tags suggestion
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.SuggestionRequest true "Content for suggestion generation"
// @Success 200 {object} SuccessResponse{data=[]models.ActivitySuggestion} "Activity suggestions retrieved successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /suggestions/activity [post]
func (h *SuggestionHandler) HandleGetActivitySuggestionsViaText(c *gin.Context) {
	var request dto.SuggestionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		FromError(c, err)
		return
	}

	suggestions, err := h.service.GetActivitySuggestionsViaText(request.Content)
	if err != nil {
		FromError(c, err)
		return
	}

	Success(c, "Activity suggestions retrieved successfully", suggestions)
}
