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
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupContributorTest(t *testing.T) (*gin.Engine, *mocks.MockContributorService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := mocks.NewMockContributorService(t)
	handler := NewContributorHandler(mockService)

	router.Use(func(c *gin.Context) {
		c.Set("claims", jwt.Claims{
			Handle: "testuser",
			Email:  "test@example.com",
		})
		c.Set("Campaign-Key", "test-key")
		c.Next()
	})

	router.POST("/contributor/:campaignID", handler.HandleAddContributor)
	router.DELETE("/contributor/:campaignID/:contributorID", handler.HandleRemoveContributor)
	router.PATCH("/contributor/:campaignID/:contributorID", handler.HandleEditContributor)
	router.GET("/contributor/:campaignID", handler.HandleGetContributorsByCampaignID)
	router.GET("/contributor/:campaignID/:contributorID", handler.HandleGetContributorByID)

	return router, mockService
}

func TestHandleAddContributor(t *testing.T) {
	router, mockService := setupContributorTest(t)

	tests := []struct {
		name           string
		campaignID     string
		contributor    models.Contributor
		setupMock      func(*mocks.MockContributorService)
		expectedCode   int
		expectedError  bool
		expectedResult string
	}{
		{
			name:       "Success",
			campaignID: "123",
			contributor: models.Contributor{
				Amount: 2000,
				Email:  "test@example.com",
				Name:   "Test Contributor",
			},
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("AddContributorToCampaign", mock.AnythingOfType("*models.Contributor"), "123", "test-key", "testuser").Return(nil)
			},
			expectedCode:   http.StatusOK,
			expectedError:  false,
			expectedResult: "Contributor added to Campaign",
		},
		{
			name:       "Service Error",
			campaignID: "123",
			contributor: models.Contributor{
				Name:   "Test Contributor",
				Amount: 2000,
				Email:  "test@example.com",
			},
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("AddContributorToCampaign", mock.AnythingOfType("*models.Contributor"), "123", "test-key", "testuser").
					Return(errors.New("service error"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedError:  true,
			expectedResult: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.setupMock(mockService)

			body, _ := json.Marshal(tt.contributor)
			req := httptest.NewRequest("POST", "/contributor/"+tt.campaignID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Equal(t, tt.expectedResult, response["message"])
			} else {
				assert.Equal(t, tt.expectedResult, response["message"])
			}
		})
	}
}

func TestHandleGetContributorByID(t *testing.T) {
	router, mockService := setupContributorTest(t)

	tests := []struct {
		name           string
		contributorID  string
		campaignID     string
		setupMock      func(*mocks.MockContributorService)
		expectedCode   int
		expectedError  bool
		expectedResult string
	}{
		{
			name:          "Success",
			contributorID: "1",
			campaignID:    "123",
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("GetContributorByID", uint(1)).Return(models.Contributor{
					Name: "Test Contributor",
				}, nil)
			},
			expectedCode:   http.StatusOK,
			expectedError:  false,
			expectedResult: "Contributor retrieved successfully",
		},
		{
			name:          "Invalid ID",
			contributorID: "invalid",
			campaignID:    "123",
			setupMock:     func(ms *mocks.MockContributorService) {},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Not Found",
			contributorID: "1",
			campaignID:    "123",
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("GetContributorByID", uint(1)).Return(models.Contributor{}, errors.New("not found"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedError:  true,
			expectedResult: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.Calls = nil
			mockService.ExpectedCalls = nil
			tt.setupMock(mockService)

			req := httptest.NewRequest("GET", "/contributor/"+tt.campaignID+"/"+tt.contributorID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusBadRequest {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				t.Log(response)

				if tt.expectedError {
					assert.Equal(t, tt.expectedResult, response["message"])
				} else {
					assert.Equal(t, tt.expectedResult, response["message"])
				}
			}
		})
	}
}

func TestHandleGetContributorsByCampaignID(t *testing.T) {
	router, mockService := setupContributorTest(t)

	tests := []struct {
		name           string
		campaignID     string
		setupMock      func(*mocks.MockContributorService)
		expectedCode   int
		expectedError  bool
		expectedResult string
	}{
		{
			name:       "Success",
			campaignID: "123",
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("GetContributorsByCampaignID", "123").Return([]models.Contributor{
					{Name: "Contributor 1"},
					{Name: "Contributor 2"},
				}, nil)
			},
			expectedCode:   http.StatusOK,
			expectedError:  false,
			expectedResult: "Contributors retrieved successfully",
		},
		{
			name:       "Service Error",
			campaignID: "123",
			setupMock: func(ms *mocks.MockContributorService) {
				ms.On("GetContributorsByCampaignID", "123").Return(nil, errors.New("service error"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedError:  true,
			expectedResult: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.setupMock(mockService)

			req := httptest.NewRequest("GET", "/contributor/"+tt.campaignID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Equal(t, tt.expectedResult, response["message"])
			} else {
				assert.Equal(t, tt.expectedResult, response["message"])
			}
		})
	}
}
