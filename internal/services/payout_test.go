package services

import (
	"errors"
	"testing"
	"time"

	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	"github.com/oyen-bright/goFundIt/internal/models"
	mockInterfaces "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	serviceMocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	loggerMock "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	paystackMock "github.com/oyen-bright/goFundIt/pkg/paystack/mocks"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPayoutService(t *testing.T) (
	*payoutService,
	*mockInterfaces.MockPayoutRepository,
	*serviceMocks.MockCampaignService,
	*serviceMocks.MockNotificationService,
	*paystackMock.MockPaystackClient,
	*serviceMocks.MockEventBroadcaster,
	*loggerMock.MockLogger,
) {
	mockRepo := mockInterfaces.NewMockPayoutRepository(t)
	mockCampaignService := serviceMocks.NewMockCampaignService(t)
	mockNotificationService := serviceMocks.NewMockNotificationService(t)
	mockPaystack := paystackMock.NewMockPaystackClient(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMock.NewMockLogger(t)

	service := NewPayoutService(
		mockRepo,
		mockCampaignService,
		mockNotificationService,
		mockPaystack,
		mockBroadcaster,
		mockLogger,
	)

	return service.(*payoutService),
		mockRepo,
		mockCampaignService,
		mockNotificationService,
		mockPaystack,
		mockBroadcaster,
		mockLogger
}

func TestInitializeManualPayout(t *testing.T) {
	service, mockRepo, mockCampaignService, mockNotificationService, _, mockBroadcaster, _ := setupPayoutService(t)

	tests := []struct {
		name        string
		campaignID  string
		userHandle  string
		setupMocks  func() *models.Campaign
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "Success",
			campaignID: "campaign1",
			userHandle: "user1",
			setupMocks: func() *models.Campaign {
				currency := models.NGN
				campaign := &models.Campaign{
					ID:           "campaign1",
					CreatedBy:    models.User{Handle: "user1"},
					FiatCurrency: &currency,
					Contributors: []models.Contributor{},
				}

				mockCampaignService.On("GetCampaignByIDWithContributors", "campaign1").Return(campaign, nil)
				mockRepo.On("Create", mock.AnythingOfType("*models.Payout")).Run(func(args mock.Arguments) {
					payout := args.Get(0).(*models.Payout)
					campaign.Payout = payout
				}).Return(nil)

				mockBroadcaster.On("NewEvent", "campaign1", websocket.EventTypePayoutUpdated, mock.AnythingOfType("*models.Payout")).Return().Once()
				mockNotificationService.On("NotifyPayoutCollected", campaign).Return(nil).Once()

				return campaign
			},
			wantErr: false,
		},
		{
			name:       "Unauthorized User",
			campaignID: "campaign1",
			userHandle: "user2",
			setupMocks: func() *models.Campaign {
				campaign := &models.Campaign{
					ID:        "campaign1",
					CreatedBy: models.User{Handle: "user1"},
				}
				mockCampaignService.On("GetCampaignByIDWithContributors", "campaign1").Return(campaign, nil)
				return campaign
			},
			wantErr:     true,
			expectedErr: "You are not authorized to perform this action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any previous mock expectations
			mockRepo.ExpectedCalls = nil
			mockCampaignService.ExpectedCalls = nil
			mockNotificationService.ExpectedCalls = nil
			mockBroadcaster.ExpectedCalls = nil

			// Setup mocks and get campaign instance
			campaign := tt.setupMocks()

			payout, err := service.InitializeManualPayout(tt.campaignID, tt.userHandle)

			// Small delay to allow goroutines to complete
			if !tt.wantErr {
				time.Sleep(100 * time.Millisecond)
			}

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, payout)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payout)
				assert.Equal(t, campaign.Payout, payout)

				// Assert all expectations were met
				mockRepo.AssertExpectations(t)
				mockCampaignService.AssertExpectations(t)
				mockNotificationService.AssertExpectations(t)
				mockBroadcaster.AssertExpectations(t)
			}
		})
	}
}

func TestVerifyAccount(t *testing.T) {
	service, _, _, _, mockPaystack, _, _ := setupPayoutService(t)

	tests := []struct {
		name        string
		request     dto.VerifyAccountRequest
		setupMocks  func()
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Success",
			request: dto.VerifyAccountRequest{
				AccountNumber: "1234567890",
				BankCode:      "001",
			},
			setupMocks: func() {
				response := &paystack.ResolveAccountResponse{
					Status: true,
					Data: struct {
						AccountNumber string "json:\"account_number\""
						AccountName   string "json:\"account_name\""
					}{
						AccountNumber: "1234567890",
						AccountName:   "Test Account",
					},
				}
				mockPaystack.EXPECT().ResolveAccount("1234567890", "001").Return(response, nil)
			},
			wantErr: false,
		},
		{
			name: "Failed Verification",
			request: dto.VerifyAccountRequest{
				AccountNumber: "1234567891",
				BankCode:      "001",
			},
			setupMocks: func() {
				response := &paystack.ResolveAccountResponse{
					Status: false,
				}
				mockPaystack.EXPECT().ResolveAccount("1234567891", "001").Return(response, nil)
			},
			wantErr:     true,
			expectedErr: "Account verification failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := service.VerifyAccount(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetBankList(t *testing.T) {
	service, _, _, _, mockPaystack, _, mockLogger := setupPayoutService(t)

	tests := []struct {
		name        string
		setupMocks  func()
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Success",
			setupMocks: func() {
				response := &paystack.BankListResponse{
					Status: true,
					Data: []paystack.Bank{
						{Code: "001", Name: "Bank 1"},
						{Code: "002", Name: "Bank 2"},
					},
				}
				mockPaystack.EXPECT().GetBanks().Return(response, nil)
			},
			wantErr: false,
		},
		{
			name: "API Error",
			setupMocks: func() {
				apiErr := errors.New("API error")
				mockPaystack.EXPECT().GetBanks().Return(nil, apiErr)
				// Mock the logger Error call
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr:     true,
			expectedErr: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous mock expectations
			mockPaystack.ExpectedCalls = nil
			mockPaystack.Calls = nil
			mockLogger.ExpectedCalls = nil

			tt.setupMocks()

			banks, err := service.GetBankList()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, banks)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, banks)
			}
		})
	}
}
