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
	mocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentHandler_HandleInitializePayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		contributorID      string
		setupMock          func(*mocks.MockPaymentService)
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:          "Success",
			contributorID: "1",
			setupMock: func(mockService *mocks.MockPaymentService) {
				payment := &models.Payment{}
				mockService.On("InitializePayment", uint(1)).Return(payment, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Payment initialized",
		},
		{
			name:          "Invalid Contributor ID",
			contributorID: "invalid",
			setupMock: func(mockService *mocks.MockPaymentService) {
				// No mock setup needed
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid contributor ID",
		},
		{
			name:          "Service Error",
			contributorID: "1",
			setupMock: func(mockService *mocks.MockPaymentService) {
				mockService.On("InitializePayment", uint(1)).Return(nil, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockPaymentService(t)
			tt.setupMock(mockService)
			handler := NewPaymentHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "contributorID", Value: tt.contributorID}}

			handler.HandleInitializePayment(c)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, response["message"])
		})
	}
}

func TestPaymentHandler_HandleVerifyPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		reference          string
		setupMock          func(*mocks.MockPaymentService)
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:      "Success",
			reference: "ref123",
			setupMock: func(mockService *mocks.MockPaymentService) {
				mockService.On("VerifyPayment", "ref123").Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Payment verified",
		},
		{
			name:      "Verification Error",
			reference: "ref123",
			setupMock: func(mockService *mocks.MockPaymentService) {
				mockService.On("VerifyPayment", "ref123").Return(errors.New("verification failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "verification failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockPaymentService(t)
			tt.setupMock(mockService)
			handler := NewPaymentHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "reference", Value: tt.reference}}

			handler.HandleVerifyPayment(c)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, response["message"])
		})
	}
}

func TestPaymentHandler_HandlePayStackWebhook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		webhookEvent       paystack.PaystackWebhookEvent
		setupMock          func(*mocks.MockPaymentService)
		expectedStatusCode int
	}{
		{
			name: "Valid Webhook",
			webhookEvent: paystack.PaystackWebhookEvent{
				Event: paystack.EventChargeSuccess,
			},
			setupMock: func(mockService *mocks.MockPaymentService) {
				mockService.On("ProcessPaystackWebhook", mock.AnythingOfType("paystack.PaystackWebhookEvent")).Return()
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockPaymentService(t)
			tt.setupMock(mockService)
			handler := NewPaymentHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set request body
			body, _ := json.Marshal(tt.webhookEvent)
			c.Request = httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.HandlePayStackWebhook(c)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
