package services

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	mockRepos "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	mockServices "github.com/oyen-bright/goFundIt/internal/services/mocks"
	loggerMock "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/oyen-bright/goFundIt/pkg/paystack"
	paystackMock "github.com/oyen-bright/goFundIt/pkg/paystack/mocks"
	storageMock "github.com/oyen-bright/goFundIt/pkg/storage/mocks"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeManualPayment(t *testing.T) {
	// Setup mocks
	mockRepo := mockRepos.NewMockPaymentRepository(t)
	mockContribService := mockServices.NewMockContributorService(t)
	mockCampaignService := mockServices.NewMockCampaignService(t)
	mockLogger := loggerMock.NewMockLogger(t)
	mockStorage := storageMock.NewMockStorage(t)
	mockNotificationService := mockServices.NewMockNotificationService(t)
	mockAnalytics := mockServices.NewMockAnalyticsService(t)
	mockBroadcaster := mockServices.NewMockEventBroadcaster(t)

	tests := []struct {
		name          string
		contributorID uint
		reference     string
		userEmail     string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:          "Valid manual payment initialization",
			contributorID: 1,
			reference:     "ref123",
			userEmail:     "user@test.com",
			setupMocks: func() {
				// Mock contributor service
				contributor := models.Contributor{
					ID:         1,
					Email:      "user@test.com",
					CampaignID: "campaign1",
					Amount:     100,
				}
				mockContribService.On("GetContributorByID", uint(1)).Return(contributor, nil)

				mockStorage.On("UploadFile", "ref123", "payment/reference").Return("url", "id", nil)

				mockRepo.On("Create", mock.AnythingOfType("*models.Payment")).Return(nil)

				mockBroadcaster.On("NewEvent",
					"campaign1",
					websocket.EventTypeContributorUpdated,
					mock.AnythingOfType("models.Contributor"),
				).Return()
			},
			expectedError: false,
		},
		{
			name:          "Already paid contributor",
			contributorID: 2,
			reference:     "ref456",
			userEmail:     "user2@test.com",
			setupMocks: func() {
				contributor := models.Contributor{
					ID:         2,
					Email:      "user2@test.com",
					CampaignID: "campaign2",
					Payment: &models.Payment{
						PaymentStatus: models.PaymentStatusSucceeded,
					}, // Has payment
				}
				mockContribService.On("GetContributorByID", uint(2)).Return(contributor, nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil
		mockCampaignService.ExpectedCalls = nil
		mockContribService.ExpectedCalls = nil
		mockContribService.Calls = nil
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMocks()

			svc := &paymentService{
				repo:                mockRepo,
				contributorService:  mockContribService,
				analyticsService:    mockAnalytics,
				campaignService:     mockCampaignService,
				notificationService: mockNotificationService,
				storage:             mockStorage,
				paystack:            nil,
				broadcaster:         mockBroadcaster,
				logger:              mockLogger,
				runAsync:            func(f func()) { f() },
			}

			// Execute test
			payment, err := svc.InitializeManualPayment(tt.contributorID, tt.reference, tt.userEmail)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)

				assert.Equal(t, tt.contributorID, payment.ContributorID)
				assert.Equal(t, models.PaymentMethodManual, payment.PaymentMethod)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockContribService.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
			mockBroadcaster.AssertExpectations(t)
		})
	}
}

func TestVerifyPayment(t *testing.T) {
	mockPaystack := paystackMock.NewMockPaystackClient(t)

	mockRepo := mockRepos.NewMockPaymentRepository(t)
	mockContribService := mockServices.NewMockContributorService(t)
	mockCampaignService := mockServices.NewMockCampaignService(t)
	mockLogger := loggerMock.NewMockLogger(t)
	mockStorage := storageMock.NewMockStorage(t)
	mockNotificationService := mockServices.NewMockNotificationService(t)
	mockAnalytics := mockServices.NewMockAnalyticsService(t)
	mockBroadcaster := mockServices.NewMockEventBroadcaster(t)

	tests := []struct {
		name          string
		reference     string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:      "Successful payment verification",
			reference: "ref123",
			setupMocks: func() {
				payment := &models.Payment{Reference: "ref123", PaymentStatus: models.PaymentStatusPending, Contributor: models.Contributor{
					CampaignID: "123",
					Payment:    &models.Payment{},
				}}
				mockRepo.On("GetByReference", "ref123").Return(payment, nil)
				mockPaystack.On("VerifyTransaction", "ref123").Return(&paystack.VerifyTransactionResponse{
					Status: true,
					Data: struct {
						ID              int    "json:\"id\""
						Status          string "json:\"status\""
						Message         string "json:\"message\""
						GatewayResponse string "json:\"gateway_response\""
						Fees            int    "json:\"fees\""
						Reference       string "json:\"reference\""
						Amount          int    "json:\"amount\""
						Currency        string "json:\"currency\""
						PaidAt          string "json:\"paid_at\""
						CreatedAt       string "json:\"created_at\""
					}{Status: "success", GatewayResponse: "Successful"},
				}, nil)
				mockRepo.On("Update", mock.AnythingOfType("*models.Payment")).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			svc := NewPaymentService(
				mockRepo,
				mockContribService,
				mockAnalytics,
				mockCampaignService,
				mockNotificationService,
				mockPaystack,
				mockStorage,
				mockBroadcaster,
				mockLogger,
			)

			err := svc.VerifyPayment(tt.reference)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyManualPayment(t *testing.T) {
	mockPaystack := paystackMock.NewMockPaystackClient(t)
	mockRepo := mockRepos.NewMockPaymentRepository(t)
	mockContribService := mockServices.NewMockContributorService(t)
	mockCampaignService := mockServices.NewMockCampaignService(t)
	mockLogger := loggerMock.NewMockLogger(t)
	mockStorage := storageMock.NewMockStorage(t)
	mockNotificationService := mockServices.NewMockNotificationService(t)
	mockAnalytics := mockServices.NewMockAnalyticsService(t)
	mockBroadcaster := mockServices.NewMockEventBroadcaster(t)

	tests := []struct {
		name          string
		reference     string
		userHandle    string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:       "Successful manual payment verification",
			reference:  "ref123",
			userHandle: "user1",
			setupMocks: func() {

				fiatCurrency := models.NGN
				campaign := &models.Campaign{
					CreatedBy:    models.User{Handle: "user1"},
					FiatCurrency: &fiatCurrency,
				}

				payment := &models.Payment{
					Reference:     "ref123",
					CampaignID:    "campaign1",
					PaymentMethod: "manual",
					Amount:        100,
					Campaign:      *campaign,
					Contributor: models.Contributor{
						CampaignID: "campaign1",
					},
				}

				// Reset mock expectations
				mockRepo.ExpectedCalls = nil
				mockCampaignService.ExpectedCalls = nil
				mockBroadcaster.ExpectedCalls = nil
				mockNotificationService.ExpectedCalls = nil
				mockAnalytics.ExpectedCalls = nil

				// Set up mock expectations
				mockRepo.On("GetByReference", "ref123").Return(payment, nil)
				mockCampaignService.On("GetCampaignByID", "campaign1").Return(campaign, nil)
				mockRepo.On("Update", mock.AnythingOfType("*models.Payment")).Return(nil)

				// Mock broadcaster with exact campaign ID
				mockBroadcaster.On("NewEvent",
					"campaign1",
					websocket.EventTypeContributorUpdated,
					mock.MatchedBy(func(c interface{}) bool {
						contrib, ok := c.(models.Contributor)
						return ok && contrib.CampaignID == "campaign1"
					}),
				).Return()

				// Mock notification service
				mockNotificationService.On("NotifyPaymentReceived",
					mock.AnythingOfType("*models.Contributor"),
					mock.AnythingOfType("*models.Campaign"),
				).Return(nil)

				mockAnalytics.On("GetCurrentData").Return(&models.PlatformAnalytics{})

			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			svc := &paymentService{
				repo:                mockRepo,
				contributorService:  mockContribService,
				analyticsService:    mockAnalytics,
				campaignService:     mockCampaignService,
				notificationService: mockNotificationService,
				storage:             mockStorage,
				paystack:            mockPaystack,
				broadcaster:         mockBroadcaster,
				logger:              mockLogger,
				runAsync:            func(f func()) { f() },
			}

			err := svc.VerifyManualPayment(tt.reference, tt.userHandle)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mock expectations
			mockRepo.AssertExpectations(t)
			mockCampaignService.AssertExpectations(t)
			mockBroadcaster.AssertExpectations(t)
			mockNotificationService.AssertExpectations(t)
		})
	}
}

func TestInitializePayment(t *testing.T) {

	mockPaystack := paystackMock.NewMockPaystackClient(t)
	mockRepo := mockRepos.NewMockPaymentRepository(t)
	mockContribService := mockServices.NewMockContributorService(t)
	mockCampaignService := mockServices.NewMockCampaignService(t)
	mockLogger := loggerMock.NewMockLogger(t)
	mockStorage := storageMock.NewMockStorage(t)
	mockNotificationService := mockServices.NewMockNotificationService(t)
	mockAnalytics := mockServices.NewMockAnalyticsService(t)
	mockBroadcaster := mockServices.NewMockEventBroadcaster(t)

	tests := []struct {
		name          string
		contributorID uint
		setupMocks    func()
		expectedError bool
	}{
		{
			name:          "Successful payment initialization",
			contributorID: 1,
			setupMocks: func() {
				fiatCurrency := models.NGN
				contributor := models.Contributor{
					ID:         1,
					Email:      "contributor-email",
					CampaignID: "campaign1",
				}
				mockContribService.On("GetContributorByID", uint(1)).Return(contributor, nil)
				campaign := &models.Campaign{
					ID:            "campaign1",
					FiatCurrency:  &fiatCurrency,
					EndDate:       time.Now().Add(24 * time.Hour),
					PaymentMethod: models.PaymentMethodFiat,
				}
				mockCampaignService.On("GetCampaignByID", "campaign1").Return(campaign, nil)
				mockPaystack.On("InitiateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&paystack.TransactionResponse{
					Data: struct {
						AuthorizationURL string "json:\"authorization_url\""
						AccessCode       string "json:\"access_code\""
						Reference        string "json:\"reference\""
					}{
						Reference:        "test-response",
						AuthorizationURL: "authrization-URL",
					},
				}, nil)
				mockRepo.On("Create", mock.AnythingOfType("*models.Payment")).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			svc := NewPaymentService(
				mockRepo,
				mockContribService,
				mockAnalytics,
				mockCampaignService,
				mockNotificationService,
				mockPaystack,
				mockStorage,
				mockBroadcaster,
				mockLogger,
			)

			payment, err := svc.InitializePayment(tt.contributorID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
			}
		})
	}
}
