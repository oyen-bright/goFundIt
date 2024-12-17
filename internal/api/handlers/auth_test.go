// internal/api/handlers/auth_test.go

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*mocks.AuthService, *handlers.AuthHandler, *gin.Engine) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup mock service
	mockService := mocks.NewAuthService(t)

	// Create handler with mock service
	handler := handlers.NewAuthHandler(mockService)

	// Setup router
	router := gin.New()

	return mockService, handler, router
}

func TestAuthHandler_HandleAuth(t *testing.T) {
	mockService, handler, router := setupTest(t)

	// Define auth endpoint
	router.POST("/auth", handler.HandleAuth)

	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "successful auth request",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
			mockSetup: func() {
				mockService.On("RequestAuth", "test@example.com", "Test User").
					Return(models.Otp{
						Email:     "test@example.com",
						RequestId: "test-request-id",
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedMsg:  "Please check your email for the OTP.",
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
				"name":  "Test User",
			},
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid inputs, please check and try again",
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"email": "error@example.com",
				"name":  "Error User",
			},
			mockSetup: func() {
				mockService.On("RequestAuth", "error@example.com", "Error User").
					Return(models.Otp{}, assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Create request body
			jsonBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			// Parse response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert message
			assert.Contains(t, response["message"], tt.expectedMsg)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_HandleVerifyAuth(t *testing.T) {
	mockService, handler, router := setupTest(t)

	// Define verify auth endpoint
	router.POST("/verify", handler.HandleVerifyAuth)

	tests := []struct {
		name          string
		requestBody   map[string]interface{}
		mockSetup     func()
		expectedCode  int
		expectedMsg   string
		expectedToken string
	}{
		{
			name: "successful verification",
			requestBody: map[string]interface{}{
				"email":      "test@example.com",
				"code":       "123456",
				"request_id": "test-request-id",
			},
			mockSetup: func() {
				mockService.On("VerifyAuth",
					"test@example.com",
					"123456",
					"test-request-id",
				).Return("valid-token", nil)
			},
			expectedCode:  http.StatusOK,
			expectedMsg:   "Authenticated",
			expectedToken: "valid-token",
		},
		{
			name: "invalid otp code",
			requestBody: map[string]interface{}{
				"email":      "test@example.com",
				"code":       "invalid",
				"request_id": "test-request-id",
			},
			mockSetup: func() {
				mockService.On("VerifyAuth",
					"test@example.com",
					"invalid",
					"test-request-id",
				).Return("", assert.AnError)
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid OTP",
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid inputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Create request body
			jsonBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			// Parse response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert message
			assert.Contains(t, response["message"], tt.expectedMsg)

			// Check token if expected
			if tt.expectedToken != "" {
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.expectedToken, data["token"])
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

// Helper function to create test request
func createTestRequest(method, url string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return httptest.NewRecorder(), req
}
