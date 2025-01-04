package services

import (
	"errors"
	"testing"

	mockAI "github.com/oyen-bright/goFundIt/internal/ai/mocks"
	"github.com/oyen-bright/goFundIt/internal/models"
	mockServices "github.com/oyen-bright/goFundIt/internal/services/mocks"
	mockLogger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSuggestionTest(t *testing.T) (*mockAI.MockAIService, *mockServices.MockCampaignService, *mockLogger.MockLogger, *suggestionService) {
	mockAIService := mockAI.NewMockAIService(t)
	mockCampaignService := mockServices.NewMockCampaignService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	service := NewSuggestionService(mockAIService, mockCampaignService, mockLogger)
	suggestionSvc, ok := service.(*suggestionService)
	if !ok {
		t.Fatal("could not cast to suggestionService")
	}

	return mockAIService, mockCampaignService, mockLogger, suggestionSvc
}

func TestGetActivitySuggestions(t *testing.T) {
	mockAI, mockCampaign, mockLogger, service := setupSuggestionTest(t)

	t.Run("Success", func(t *testing.T) {
		campaignID := "test-id"
		campaignKey := "test-key"
		campaign := &models.Campaign{
			ID:          campaignID,
			Title:       "Test Campaign",
			Description: "Test Description",
		}
		expectedSuggestions := []models.ActivitySuggestion{
			{Title: "Suggestion 1"},
			{Title: "Suggestion 2"},
		}

		mockCampaign.EXPECT().GetCampaignByID(campaignID, campaignKey).Return(campaign, nil)
		mockAI.EXPECT().GenerateActivitySuggestions(campaign.Title+" "+campaign.Description).
			Return(expectedSuggestions, nil)

		suggestions, err := service.GetActivitySuggestions(campaignID, campaignKey)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuggestions, suggestions)
	})

	t.Run("Campaign Not Found", func(t *testing.T) {
		mockAI.ExpectedCalls = nil
		mockAI.Calls = nil
		mockCampaign.ExpectedCalls = nil
		mockCampaign.Calls = nil
		mockLogger.ExpectedCalls = nil
		mockLogger.Calls = nil
		campaignID := "invalid-id"
		campaignKey := "test-key"
		mockCampaign.EXPECT().GetCampaignByID(campaignID, campaignKey).Return(nil, errors.New("campaign not found"))

		suggestions, err := service.GetActivitySuggestions(campaignID, campaignKey)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})

	t.Run("AI Service Error", func(t *testing.T) {
		mockAI.ExpectedCalls = nil
		mockAI.Calls = nil
		mockCampaign.ExpectedCalls = nil
		mockCampaign.Calls = nil
		mockLogger.ExpectedCalls = nil
		mockLogger.Calls = nil
		campaignID := "test-id"
		campaignKey := "test-key"
		campaign := &models.Campaign{
			ID:          campaignID,
			Title:       "Test Campaign",
			Description: "Test Description",
		}

		mockCampaign.EXPECT().GetCampaignByID(campaignID, campaignKey).Return(campaign, nil)
		mockAI.EXPECT().GenerateActivitySuggestions(campaign.Title+" "+campaign.Description).
			Return(nil, errors.New("ai service error"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything)

		suggestions, err := service.GetActivitySuggestions(campaignID, campaignKey)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})
}

func TestGetActivitySuggestionsViaText(t *testing.T) {
	mockAI, _, mockLogger, service := setupSuggestionTest(t)

	t.Run("Success", func(t *testing.T) {
		content := "test content"
		expectedSuggestions := []models.ActivitySuggestion{
			{Title: "Suggestion 1"},
			{Title: "Suggestion 2"},
		}

		mockAI.EXPECT().GenerateActivitySuggestions(content).Return(expectedSuggestions, nil)

		suggestions, err := service.GetActivitySuggestionsViaText(content)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuggestions, suggestions)
	})

	t.Run("AI Service Error", func(t *testing.T) {
		mockAI.ExpectedCalls = nil
		mockAI.Calls = nil
		mockLogger.ExpectedCalls = nil
		mockLogger.Calls = nil

		content := "test content"
		mockAI.EXPECT().GenerateActivitySuggestions(content).Return(nil, errors.New("ai service error"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything)

		suggestions, err := service.GetActivitySuggestionsViaText(content)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})
}
