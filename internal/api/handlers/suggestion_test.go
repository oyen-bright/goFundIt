package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	interfaces "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetActivitySuggestions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		campaignID     string
		setupMock      func(*interfaces.MockSuggestionService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:       "Success",
			campaignID: "123",
			setupMock: func(m *interfaces.MockSuggestionService) {
				m.On("GetActivitySuggestions", "123", "123").Return([]models.ActivitySuggestion{
					{Title: "title", EstimatedPrice: "cost"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Activity suggestions retrieved successfully",
		},
		{
			name:       "Error",
			campaignID: "123",
			setupMock: func(m *interfaces.MockSuggestionService) {
				m.On("GetActivitySuggestions", "123", "123").Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.NewMockSuggestionService(t)
			tt.setupMock(mockService)
			handler := NewSuggestionHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, gin.Param{Key: "campaignID", Value: tt.campaignID})
			c.Set("Campaign-Key", tt.campaignID)

			handler.HandleGetActivitySuggestions(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedMsg != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response["message"])
			}
		})
	}
}

func TestHandleGetActivitySuggestionsViaText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*interfaces.MockSuggestionService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:        "Success",
			requestBody: `{"content": "test content"}`,
			setupMock: func(m *interfaces.MockSuggestionService) {
				m.On("GetActivitySuggestionsViaText", "test content").Return([]models.ActivitySuggestion{
					{Title: "title", EstimatedPrice: "cost"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Activity suggestions retrieved successfully",
		},
		{
			name:        "Error",
			requestBody: `{"content": "test content"}`,
			setupMock: func(m *interfaces.MockSuggestionService) {
				m.On("GetActivitySuggestionsViaText", "test content").Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid Request",
			requestBody:    `{"invalid": "json"}`,
			setupMock:      func(m *interfaces.MockSuggestionService) {},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.NewMockSuggestionService(t)
			tt.setupMock(mockService)
			handler := NewSuggestionHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.HandleGetActivitySuggestionsViaText(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedMsg != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response["message"])
			}
		})
	}
}
