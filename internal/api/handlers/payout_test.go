package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	"github.com/oyen-bright/goFundIt/internal/models"
	mock "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PayoutHandlerTestSuite struct {
	suite.Suite
	router  *gin.Engine
	mock    *mock.MockPayoutService
	handler *PayoutHandler
}

func (suite *PayoutHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.mock = mock.NewMockPayoutService(suite.T())
	suite.handler = NewPayoutHandler(suite.mock)
}

func TestPayoutHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutHandlerTestSuite))
}

func (suite *PayoutHandlerTestSuite) TestHandleGetBankList() {
	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			setupMock: func() {
				banks := []interface{}{
					map[string]interface{}{"name": "Bank1", "code": "001"},
					map[string]interface{}{"name": "Bank2", "code": "002"},
				}
				suite.mock.EXPECT().GetBankList().Return(banks, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Bank list retrieved successfully",
				"data": []interface{}{
					map[string]interface{}{"name": "Bank1", "code": "001"},
					map[string]interface{}{"name": "Bank2", "code": "002"},
				},
			},
		},
		{
			name: "error",
			setupMock: func() {
				suite.mock.EXPECT().GetBankList().Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "service error",
				"status":  "Internal Server Error",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.mock.ExpectedCalls = nil
			suite.mock.Calls = nil
			tc.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			suite.handler.HandleGetBankList(c)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tc.expectedBody, response)
		})
	}
}

func (suite *PayoutHandlerTestSuite) TestHandleVerifyAccount() {
	testReq := dto.VerifyAccountRequest{
		AccountNumber: "1234567890",
		BankCode:      "001",
	}

	tests := []struct {
		name           string
		request        interface{}
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "success",
			request: testReq,
			setupMock: func() {
				resp := map[string]interface{}{"account_name": "Test Account"}
				suite.mock.EXPECT().VerifyAccount(testReq).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Account verified successfully",
				"data":    map[string]interface{}{"account_name": "Test Account"},
			},
		},
		{
			name:    "invalid request",
			request: map[string]interface{}{},
			setupMock: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			tc.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tc.request)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			suite.handler.HandleVerifyAccount(c)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)
			if tc.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedBody, response)
			}
		})
	}
}

func (suite *PayoutHandlerTestSuite) TestHandleInitializePayout() {
	testReq := dto.PayoutRequest{
		AccountName:   "Test Account",
		AccountNumber: "1234567890",
		BankName:      "Test Bank",
		BankCode:      "001",
	}

	tests := []struct {
		name           string
		request        interface{}
		setupContext   func(*gin.Context)
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "success",
			request: testReq,
			setupContext: func(c *gin.Context) {

				c.Set("claims", jwt.Claims{Handle: "user123"})
			},
			setupMock: func() {
				payout := &models.Payout{ID: "payout123", Amount: 0.0}
				suite.mock.EXPECT().InitializePayout("campaign123", "user123", testReq).Return(payout, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Payout initialized successfully",
				"data": map[string]interface{}{
					"amount":        0.0,
					"campaignId":    "",
					"completedAt":   interface{}(nil),
					"processedAt":   interface{}(nil),
					"reference":     "",
					"payoutMethod":  "",
					"status":        "",
					"failureReason": interface{}(nil),
				},
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			tc.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			jsonData, _ := json.Marshal(tc.request)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			c.Params = gin.Params{
				{Key: "campaignID", Value: "campaign123"},
			}
			suite.handler.HandleInitializePayout(c)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)
			if tc.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedBody, response)
			}
		})
	}
}

func (suite *PayoutHandlerTestSuite) TestHandleGetPayoutByCampaignID() {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			setupContext: func(c *gin.Context) {
			},
			setupMock: func() {
				payout := &models.Payout{ID: "payout123", Amount: 0.0}
				suite.mock.EXPECT().GetPayoutByCampaignID("campaign123").Return(payout, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Payout information retrieved successfully", "status": "OK",
				"data": map[string]interface{}{
					"amount": 0.0, "campaignId": "", "completedAt": interface{}(nil), "failureReason": interface{}(nil), "payoutMethod": "", "processedAt": interface{}(nil), "reference": "", "status": ""},
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			tc.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			c.Params = gin.Params{
				{Key: "campaignID", Value: "campaign123"},
			}

			suite.handler.HandleGetPayoutByCampaignID(c)

			assert.Equal(suite.T(), tc.expectedStatus, w.Code)
			if tc.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedBody, response)
			}
		})
	}
}
