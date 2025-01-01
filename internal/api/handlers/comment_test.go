package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	mocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*mocks.MockCommentService)
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"content": "Test comment",
			},
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().CreateComment(mock.AnythingOfType("*models.Comment"), "test-campaign", uint(1), "test-user").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]interface{}{
				"content": "",
			},
			setupMock: func(m *mocks.MockCommentService) {
				// No mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			requestBody: map[string]interface{}{
				"content": "Test comment",
			},
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().CreateComment(mock.AnythingOfType("*models.Comment"), "test-campaign", uint(1), "test-user").
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCommentService(t)
			tt.setupMock(mockService)

			handler := &CommentHandler{
				CommentService: mockService,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request
			body, _ := json.Marshal(tt.requestBody)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

			c.Set("claims", jwt.Claims{Handle: "test-user"})
			c.Set("Campaign-Key", "test-campaign")
			c.Params = []gin.Param{
				{Key: "campaignID", Value: "test-campaign"},
				{Key: "activityID", Value: "1"},
			}

			handler.HandleCreateComment(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandleGetActivityComments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockCommentService)
		expectedStatus int
	}{
		{
			name: "Success",
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().GetActivityComments(uint(1)).
					Return([]models.Comment{{Content: "Test comment"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Service Error",
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().GetActivityComments(uint(1)).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCommentService(t)
			tt.setupMock(mockService)

			handler := &CommentHandler{
				CommentService: mockService,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("claims", jwt.Claims{Handle: "test-user"})

			c.Params = []gin.Param{
				{Key: "campaignID", Value: "test-campaign"},
				{Key: "activityID", Value: "1"},
			}

			handler.HandleGetActivityComments(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandleDeleteComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockCommentService)
		expectedStatus int
	}{
		{
			name: "Success",
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().DeleteComment("1", "test-user").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Service Error",
			setupMock: func(m *mocks.MockCommentService) {
				m.EXPECT().DeleteComment("1", "test-user").
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockCommentService(t)
			mockService.Calls = nil
			mockService.ExpectedCalls = nil
			tt.setupMock(mockService)

			handler := &CommentHandler{
				CommentService: mockService,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("claims", jwt.Claims{Handle: "test-user"})
			c.Params = []gin.Param{
				{Key: "campaignID", Value: "test-campaign"},

				{Key: "commentID", Value: "1"},
			}

			handler.HandleDeleteComment(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
