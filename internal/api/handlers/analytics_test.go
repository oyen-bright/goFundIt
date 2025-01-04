package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func setupAnalyticsTest() (*mocks.MockAnalyticsService, *AnalyticsHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewMockAnalyticsService(new(testing.T))
	handler := NewAnalyticsHandler(mockService)
	router := gin.New()
	return mockService, handler, router
}

func TestHandleProcessAnalyticsNow(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockAnalyticsService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func(m *mocks.MockAnalyticsService) {
				m.EXPECT().ProcessAnalyticsNow().Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Analytics processed and sent to email","status":"OK"}`,
		},
		{
			name: "Service Error",
			setupMock: func(m *mocks.MockAnalyticsService) {
				m.EXPECT().ProcessAnalyticsNow().Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"assert.AnError general error for testing","status":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupAnalyticsTest()
			tt.setupMock(mockService)

			router.POST("/analytics/process", handler.HandleProcessAnalyticsNow)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/analytics/process", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
