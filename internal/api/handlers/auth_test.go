package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	mock_interfaces "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func setupAuthTest(t *testing.T) (*gin.Engine, *mock_interfaces.MockAuthService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := mock_interfaces.NewMockAuthService(t)
	handler := NewAuthHandler(mockService)

	router.POST("/auth", handler.HandleAuth)
	router.POST("/verify", handler.HandleVerifyAuth)
	router.POST("/fcm/save-token", func(c *gin.Context) {
		// Mock the auth middleware by setting claims
		c.Set("claims", jwt.Claims{Handle: "test-handle"})
		handler.HandleSaveFCMToken(c)
	})

	return router, mockService
}

func TestHandleAuth(t *testing.T) {
	router, mockService := setupAuthTest(t)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful auth request",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
			mockSetup: func() {
				mockService.EXPECT().RequestAuth("test@example.com", "Test User").
					Return(models.Otp{Email: "test@example.com", RequestId: "123"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Please check your email for the OTP.",
				"data": map[string]interface{}{
					"requestId": "123",
				},
			},
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
				"name":  "Test User",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
			mockSetup: func() {
				mockService.EXPECT().RequestAuth("test@example.com", "Test User").
					Return(models.Otp{}, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandleVerifyAuth(t *testing.T) {
	router, mockService := setupAuthTest(t)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful verification",
			requestBody: map[string]interface{}{
				"email":     "test@example.com",
				"code":      "123456",
				"requestId": "req123",
			},
			mockSetup: func() {
				mockService.EXPECT().VerifyAuth("test@example.com", "123456", "req123").
					Return("token123", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Authenticated",
				"data": map[string]interface{}{
					"token": "token123",
				},
			},
		},
		{
			name: "invalid verification",
			requestBody: map[string]interface{}{
				"email":     "test@example.com",
				"code":      "123456",
				"requestId": "req123",
			},
			mockSetup: func() {
				mockService.EXPECT().VerifyAuth("test@example.com", "123456", "req123").
					Return("", errors.New("invalid code"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandleSaveFCMToken(t *testing.T) {
	router, mockService := setupAuthTest(t)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful token save",
			requestBody: map[string]interface{}{
				"fcmToken": "fcm-token-123",
			},
			mockSetup: func() {
				mockService.EXPECT().SaveFCMToken("test-handle", "fcm-token-123").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "FCM token saved",
			},
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"fcmToken": "fcm-token-123",
			},
			mockSetup: func() {
				mockService.EXPECT().SaveFCMToken("test-handle", "fcm-token-123").
					Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/fcm/save-token", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
