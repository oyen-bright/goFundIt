package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type AIService interface {
	GenerateActivitySuggestions(campaignDescription string) ([]models.ActivitySuggestion, error)
}
