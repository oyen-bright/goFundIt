package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	mocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func setupActivityTest() (*mocks.MockActivityService, *ActivityHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewMockActivityService(new(testing.T))
	handler := NewActivityHandler(mockService)
	router := gin.New()
	return mockService, handler, router
}

func setupAuthMiddleware(c *gin.Context) {
	claims := jwt.Claims{
		Handle: "testuser",
		Email:  "test@example.com",
	}
	c.Set("claims", claims)
	c.Set("Campaign-Key", "test-campaign-key")
	c.Next()
}

func TestHandleGetActivitiesByCampaignID(t *testing.T) {
	tests := []struct {
		name           string
		campaignID     string
		setupMock      func(*mocks.MockActivityService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Success",
			campaignID: "campaign123",
			setupMock: func(m *mocks.MockActivityService) {
				activities := []models.Activity{{
					ID:           1,
					Title:        "Test Activity",
					CampaignID:   "",
					ImageUrl:     "",
					Subtitle:     "",
					Cost:         0,
					IsMandatory:  false,
					IsApproved:   false,
					Contributors: nil,
				}}
				m.EXPECT().GetActivitiesByCampaignID("campaign123").Return(activities, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"data": [{
					"id": 1,
					"title": "Test Activity",
					"campaignId": "",
					"imageUrl": "",
					"subtitle": "",
					"cost": 0,
					"isMandatory": false,
					"isApproved": false,
					"contributors": null
				}],
				"message": "Activities fetched successfully",
				"status": "OK"
			}`,
		},
		{
			name:       "Service Error",
			campaignID: "campaign123",
			setupMock: func(m *mocks.MockActivityService) {
				m.EXPECT().GetActivitiesByCampaignID("campaign123").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"assert.AnError general error for testing","status":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupActivityTest()
			tt.setupMock(mockService)

			router.GET("/activities/:campaignID", func(c *gin.Context) {
				setupAuthMiddleware(c)
				handler.HandleGetActivitiesByCampaignID(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/activities/%s", tt.campaignID), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandleCreateActivity(t *testing.T) {
	tests := []struct {
		name           string
		activity       models.Activity
		setupMock      func(*mocks.MockActivityService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			activity: models.Activity{
				Title: "New Activity",
				Cost:  100,
			},
			setupMock: func(m *mocks.MockActivityService) {
				expectedActivity := models.Activity{
					Title: "New Activity",
					Cost:  100,
				}
				returnedActivity := models.Activity{
					ID:           1,
					Title:        "New Activity",
					Cost:         100,
					CampaignID:   "",
					ImageUrl:     "",
					Subtitle:     "",
					IsMandatory:  false,
					IsApproved:   false,
					Contributors: nil,
				}
				m.EXPECT().CreateActivity(expectedActivity, "testuser", "campaign123", "test-campaign-key").
					Return(returnedActivity, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"data": {
					"id": 1,
					"title": "New Activity",
					"cost": 100,
					"campaignId": "",
					"imageUrl": "",
					"subtitle": "",
					"isMandatory": false,
					"isApproved": false,
					"contributors": null
				},
				"message": "Activity created successfully",
				"status": "OK"
			}`,
		},
		{
			name: "Service Error",
			activity: models.Activity{
				Title: "New Activity",
				Cost:  100,
			},
			setupMock: func(m *mocks.MockActivityService) {
				expectedActivity := models.Activity{
					Title: "New Activity",
					Cost:  100,
				}
				m.EXPECT().CreateActivity(expectedActivity, "testuser", "campaign123", "test-campaign-key").
					Return(models.Activity{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"assert.AnError general error for testing","status":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupActivityTest()
			tt.setupMock(mockService)

			router.POST("/activities/:campaignID", func(c *gin.Context) {
				setupAuthMiddleware(c)
				handler.HandleCreateActivity(c)
			})

			activityJSON, _ := json.Marshal(tt.activity)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/activities/campaign123", bytes.NewBuffer(activityJSON))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandleGetActivityByID(t *testing.T) {
	tests := []struct {
		name           string
		activityID     uint
		campaignID     string
		setupMock      func(*mocks.MockActivityService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Success",
			activityID: 1,
			campaignID: "campaign123",
			setupMock: func(m *mocks.MockActivityService) {
				activity := models.Activity{
					ID:           1,
					Title:        "Test Activity",
					CampaignID:   "",
					ImageUrl:     "",
					Subtitle:     "",
					Cost:         0,
					IsMandatory:  false,
					IsApproved:   false,
					Contributors: nil,
				}
				m.EXPECT().GetActivityByID(uint(1), "campaign123").Return(activity, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"data": {
					"id": 1,
					"title": "Test Activity",
					"campaignId": "",
					"imageUrl": "",
					"subtitle": "",
					"cost": 0,
					"isMandatory": false,
					"isApproved": false,
					"contributors": null
				},
				"message": "Activity fetched successfully",
				"status": "OK"
			}`,
		},
		{
			name:       "Not Found",
			activityID: 999,
			campaignID: "campaign123",
			setupMock: func(m *mocks.MockActivityService) {
				m.EXPECT().GetActivityByID(uint(999), "campaign123").Return(models.Activity{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"assert.AnError general error for testing","status":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupActivityTest()
			tt.setupMock(mockService)

			router.GET("/activities/:campaignID/:activityID", func(c *gin.Context) {
				setupAuthMiddleware(c)
				handler.HandleGetActivityByID(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/activities/%s/%d", tt.campaignID, tt.activityID), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
