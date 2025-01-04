package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/campaign"
	"github.com/oyen-bright/goFundIt/internal/models"
	mocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateCampaign(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*mocks.MockCampaignService)
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"title":       "Trip to New York",
				"description": "Exciting journey to explore the Big Apple!, Exciting journey to explore the Big Apple, Exciting journey to explore the Big Apple,Exciting journey to explore the Big Apple",
				"images": []map[string]string{
					{
						"ImageUrl": "https://example.com/newyork.jpg",
					},
				},
				"PaymentMethod": "manual",
				"fiatCurrency":  "NGN",
				"Activities": []map[string]interface{}{
					{
						"title":       "Times Square Visit",
						"cost":        100,
						"IsMandatory": true,
						"isApproved":  true,
					},
				},
				"Contributors": []map[string]interface{}{
					{
						"amount": 2000,
						"Email":  "bright@krotrust.com",
					},
				},
				"StartDate": "2025-01-01T00:00:00Z",
				"EndDate":   "2025-01-03T00:00:00Z",
			},
			setupMock: func(m *mocks.MockCampaignService) {
				m.EXPECT().CreateCampaign(mock.AnythingOfType("*models.Campaign"), "test-user").
					Return(models.Campaign{Title: "test"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]interface{}{
				"title": "",
			},
			setupMock: func(m *mocks.MockCampaignService) {
				// No mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		// {
		// 	name: "Service Error",
		// 	requestBody: map[string]interface{}{
		// 		"title":       "Test Campaign",
		// 		"description": "Test Description",
		// 	},
		// 	setupMock: func(m *mocks.MockCampaignService) {
		// 		m.EXPECT().CreateCampaign(mock.AnythingOfType("*models.Campaign"), "test-user").
		// 			Return(models.Campaign{}, errors.New("service error"))
		// 	},
		// 	expectedStatus: http.StatusBadRequest,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCampaignService(t)
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.setupMock(mockService)

			handler := NewCampaignHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set claims in context
			c.Set("claims", jwt.Claims{Handle: "test-user"})
			c.Next()

			// Call the handler
			handler.HandleCreateCampaign(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// For successful cases, verify response structure
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check if the response contains the expected fields
				assert.NotNil(t, response["message"])
				assert.NotNil(t, response["data"])
				assert.Equal(t, "Campaign created successfully", response["message"])
			}
		})
	}
}

func TestHandleGetCampaignByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		campaignID     string
		setupMock      func(*mocks.MockCampaignService)
		expectedStatus int
	}{
		{
			name:       "Success",
			campaignID: "test-campaign",
			setupMock: func(m *mocks.MockCampaignService) {
				m.EXPECT().GetCampaignByID("test-campaign", "test-key").
					Return(&models.Campaign{Title: "Test Campaign"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Campaign Not Found",
			campaignID: "non-existent",
			setupMock: func(m *mocks.MockCampaignService) {
				m.EXPECT().GetCampaignByID("non-existent", "test-key").
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCampaignService(t)
			tt.setupMock(mockService)

			handler := NewCampaignHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			c.Params = []gin.Param{{Key: "campaignID", Value: tt.campaignID}}
			c.Set("Campaign-Key", "test-key")
			c.Next()

			handler.HandleGetCampaignByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandleUpdateCampaignByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	title := "Updated Title"
	updateRequest := dto.CampaignUpdateRequest{
		Title: &title,
	}

	tests := []struct {
		name           string
		campaignID     string
		requestBody    interface{}
		setupMock      func(*mocks.MockCampaignService)
		expectedStatus int
	}{
		{
			name:        "Success",
			campaignID:  "test-campaign",
			requestBody: updateRequest,
			setupMock: func(m *mocks.MockCampaignService) {
				m.EXPECT().UpdateCampaign(mock.AnythingOfType("dto.CampaignUpdateRequest"), "test-campaign", "test-user", "test-key").
					Return(&models.Campaign{Title: "Updated Title"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid Request Body",
			campaignID:  "test-campaign",
			requestBody: "invalid",
			setupMock: func(m *mocks.MockCampaignService) {
				// No mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Service Error",
			campaignID:  "test-campaign",
			requestBody: updateRequest,
			setupMock: func(m *mocks.MockCampaignService) {
				m.EXPECT().UpdateCampaign(mock.AnythingOfType("dto.CampaignUpdateRequest"), "test-campaign", "test-user", "test-key").
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCampaignService(t)
			tt.setupMock(mockService)

			handler := NewCampaignHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(body))
			c.Params = []gin.Param{{Key: "campaignID", Value: tt.campaignID}}

			// Mock JWT claims
			c.Set("claims", jwt.Claims{Handle: "test-user"})
			c.Set("Campaign-Key", "test-key")

			handler.HandleUpdateCampaignByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
