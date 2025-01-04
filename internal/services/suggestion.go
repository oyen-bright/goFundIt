package services

import (
	ai "github.com/oyen-bright/goFundIt/internal/ai/interfaces"
	"github.com/oyen-bright/goFundIt/internal/models"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type suggestionService struct {
	aIService       ai.AIService
	campaignService services.CampaignService
	logger          logger.Logger
}

// NewSuggestionService creates a new instance of the suggestion service
func NewSuggestionService(
	aIService ai.AIService,
	campaignService services.CampaignService,
	logger logger.Logger,
) services.SuggestionService {
	return &suggestionService{
		aIService:       aIService,
		campaignService: campaignService,
		logger:          logger,
	}
}

// GenerateActivitySuggestions generates activity suggestions based on the campaign description.
func (s *suggestionService) GetActivitySuggestions(campaignID, key string) ([]models.ActivitySuggestion, error) {
	//Validate campaign

	campaign, err := s.campaignService.GetCampaignByID(campaignID, key)
	if err != nil {
		return nil, err
	}

	//Get suggestions from AI service
	suggestions, err := s.aIService.GenerateActivitySuggestions(campaign.Title + " " + campaign.Description)

	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}

	return suggestions, nil

}

func (s *suggestionService) GetActivitySuggestionsViaText(content string) ([]models.ActivitySuggestion, error) {
	suggestions, err := s.aIService.GenerateActivitySuggestions(content)

	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}

	return suggestions, nil
}
