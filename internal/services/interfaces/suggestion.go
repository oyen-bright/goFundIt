package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type SuggestionService interface {
	GetActivitySuggestions(campaignID, key string) ([]models.ActivitySuggestion, error)
	GetActivitySuggestionsViaText(content string) ([]models.ActivitySuggestion, error)
}
