package services

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	mockRepo "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	mockInterfaces "github.com/oyen-bright/goFundIt/internal/services/mocks"
	mockLogger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"

	"github.com/oyen-bright/goFundIt/pkg/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateActivity(t *testing.T) {
	// Setup mocks
	mockRepo := mockRepo.NewMockActivityRepository(t)
	mockAuth := mockInterfaces.NewMockAuthService(t)
	mockCampaign := mockInterfaces.NewMockCampaignService(t)
	mockBroadcaster := mockInterfaces.NewMockEventBroadcaster(t)
	mockAnalytics := mockInterfaces.NewMockAnalyticsService(t)
	mockNotification := mockInterfaces.NewMockNotificationService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	// Create service
	service := &activityService{
		repo:                mockRepo,
		authService:         mockAuth,
		campaignService:     mockCampaign,
		broadcaster:         mockBroadcaster,
		analyticsService:    mockAnalytics,
		notificationService: mockNotification,
		logger:              mockLogger,
		runAsync:            func(f func()) { f() },
	}

	tests := []struct {
		name        string
		activity    models.Activity
		userHandle  string
		campaignID  string
		campaignKey string
		setupMocks  func()
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful creation - campaign owner",
			activity: models.Activity{
				Title:     "Test Activity",
				Cost:      100,
				CreatedBy: models.User{Handle: "user1"},
			},
			userHandle:  "user1",
			campaignKey: "campaignKey",
			campaignID:  "campaign1",
			setupMocks: func() {
				// Mock GetCampaignByID
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaignKey").Return(
					&models.Campaign{
						ID:           "campaign1",
						CreatedBy:    models.User{Handle: "user1"},
						Contributors: []models.Contributor{{Email: "user1@test.com"}},
					}, nil,
				)

				// Mock GetUserByHandle
				mockAuth.EXPECT().GetUserByHandle("user1").Return(
					models.User{
						Handle: "user1",
						Email:  "user1@test.com",
					}, nil,
				)

				// Mock Create
				mockRepo.EXPECT().Create(mock.AnythingOfType("*models.Activity")).Return(
					models.Activity{
						Title:           "Test Activity",
						Cost:            100,
						CreatedByHandle: "user1",
						CreatedBy:       models.User{Handle: "user1"},
						IsApproved:      true,
					}, nil,
				)

				// Mock broadcaster
				mockBroadcaster.EXPECT().NewEvent(
					"campaign1",
					websocket.EventTypeActivityCreated,
					mock.AnythingOfType("models.Activity"),
				)

				// Mock notifications
				mockNotification.EXPECT().NotifyActivityAddition(
					mock.AnythingOfType("*models.Activity"),
					mock.AnythingOfType("*models.Campaign"),
				).Return(nil)

				// Mock analytics
				platformAnalytics := &models.PlatformAnalytics{}
				mockAnalytics.EXPECT().GetCurrentData().Return(platformAnalytics)
			},
			wantErr: false,
		},
		{
			name: "creation fails - user not part of campaign",
			activity: models.Activity{
				Title: "Test Activity",
				Cost:  100,
			},
			userHandle:  "user2",
			campaignID:  "campaign1",
			campaignKey: "campaignKey",
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaignKey").Return(
					&models.Campaign{
						ID:           "campaign1",
						CreatedBy:    models.User{Handle: "user1"},
						Contributors: []models.Contributor{{Email: "user1@test.com"}},
					}, nil,
				)

				mockAuth.EXPECT().GetUserByHandle("user2").Return(
					models.User{
						Handle: "user2",
						Email:  "user2@test.com",
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Sorry, you can't add activities to campaigns you're not part of. Join the campaign to get started!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCampaign.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			tt.setupMocks()

			result, err := service.CreateActivity(tt.activity, tt.userHandle, tt.campaignID, tt.campaignKey)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Equal(t, tt.activity.Title, result.Title)
				assert.Equal(t, tt.userHandle, result.CreatedByHandle)
				assert.True(t, result.IsApproved)
			}
		})
	}
}

func TestApproveActivity(t *testing.T) {
	// Setup mocks
	mockRepo := mockRepo.NewMockActivityRepository(t)
	mockAuth := mockInterfaces.NewMockAuthService(t)
	mockCampaign := mockInterfaces.NewMockCampaignService(t)
	mockBroadcaster := mockInterfaces.NewMockEventBroadcaster(t)
	mockAnalytics := mockInterfaces.NewMockAnalyticsService(t)
	mockNotification := mockInterfaces.NewMockNotificationService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	// Create service
	service := NewActivityService(
		mockRepo,
		mockAuth,
		mockCampaign,
		mockBroadcaster,
		mockAnalytics,
		mockNotification,
		mockLogger,
	)

	tests := []struct {
		name        string
		activityID  uint
		userHandle  string
		setupMocks  func()
		campaignKey string
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "successful approval",
			activityID:  1,
			campaignKey: "campaign-key",
			userHandle:  "owner",
			setupMocks: func() {
				// Mock GetByID
				mockRepo.EXPECT().GetByID(uint(1)).Return(
					models.Activity{
						ID:         1,
						CampaignID: "campaign1",
						IsApproved: false,
					}, nil,
				)

				// Mock GetCampaignByID
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID:        "campaign1",
						CreatedBy: models.User{Handle: "owner"},
					}, nil,
				)

				// Mock Update
				mockRepo.EXPECT().Update(mock.AnythingOfType("*models.Activity")).Return(nil)

				// Mock broadcaster
				mockBroadcaster.EXPECT().NewEvent(
					"campaign1",
					websocket.EventTypeActivityUpdated,
					mock.AnythingOfType("models.Activity"),
				)

				// Mock notifications
				mockNotification.EXPECT().NotifyActivityApproved(
					mock.AnythingOfType("*models.Activity"),
					mock.AnythingOfType("*models.Campaign"),
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "unauthorized approval",
			activityID:  1,
			campaignKey: "campaign-key",
			userHandle:  "notowner",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(uint(1)).Return(
					models.Activity{
						ID:         1,
						CampaignID: "campaign1",
						IsApproved: false,
					}, nil,
				)

				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID:        "campaign1",
						CreatedBy: models.User{Handle: "owner"},
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Unauthorize: only campaign creator can approve activity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			activity, err := service.ApproveActivity(tt.activityID, tt.userHandle, tt.campaignKey)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, activity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, activity)
				assert.True(t, activity.IsApproved)
			}
		})
	}
}

func TestOptInContributor(t *testing.T) {
	// Setup mocks
	mockRepo := mockRepo.NewMockActivityRepository(t)
	mockAuth := mockInterfaces.NewMockAuthService(t)
	mockCampaign := mockInterfaces.NewMockCampaignService(t)
	mockBroadcaster := mockInterfaces.NewMockEventBroadcaster(t)
	mockAnalytics := mockInterfaces.NewMockAnalyticsService(t)
	mockNotification := mockInterfaces.NewMockNotificationService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	service := &activityService{
		repo:                mockRepo,
		authService:         mockAuth,
		campaignService:     mockCampaign,
		broadcaster:         mockBroadcaster,
		analyticsService:    mockAnalytics,
		notificationService: mockNotification,
		logger:              mockLogger,
		runAsync:            func(f func()) { f() },
	}

	tests := []struct {
		name          string
		campaignID    string
		userEmail     string
		activityID    uint
		contributorID uint
		setupMocks    func()
		campaignKey   string
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "successful opt-in",
			campaignID:    "campaign1",
			userEmail:     "user@test.com",
			campaignKey:   "campaign-key",
			activityID:    1,
			contributorID: 1,
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{ID: 1, IsApproved: true},
						},
						Contributors: []models.Contributor{
							{ID: 1, Email: "user@test.com", Payment: nil},
						},
					}, nil,
				)

				mockRepo.EXPECT().AddContributor(uint(1), uint(1)).Return(nil)

				mockBroadcaster.EXPECT().NewEvent(
					"campaign1",
					websocket.EventTypeActivityUpdated,
					mock.AnythingOfType("*models.Activity"),
				)
			},
			wantErr: false,
		},
		{
			name:          "opt-in fails - already paid",
			campaignID:    "campaign1",
			userEmail:     "user@test.com",
			campaignKey:   "campaign-key",
			activityID:    1,
			contributorID: 1,
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{ID: 1, IsApproved: true},
						},
						Contributors: []models.Contributor{
							{ID: 1, Email: "user@test.com", Payment: &models.Payment{
								PaymentStatus: models.PaymentStatusSucceeded,
							}},
						},
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Action cannot be performed after making a payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCampaign.ExpectedCalls = nil
			tt.setupMocks()

			err := service.OptInContributor(tt.campaignID, tt.userEmail, tt.campaignKey, tt.activityID, tt.contributorID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOptOutContributor(t *testing.T) {
	mockRepo := mockRepo.NewMockActivityRepository(t)
	mockAuth := mockInterfaces.NewMockAuthService(t)
	mockCampaign := mockInterfaces.NewMockCampaignService(t)
	mockBroadcaster := mockInterfaces.NewMockEventBroadcaster(t)
	mockAnalytics := mockInterfaces.NewMockAnalyticsService(t)
	mockNotification := mockInterfaces.NewMockNotificationService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	service := &activityService{
		repo:                mockRepo,
		authService:         mockAuth,
		campaignService:     mockCampaign,
		broadcaster:         mockBroadcaster,
		analyticsService:    mockAnalytics,
		notificationService: mockNotification,
		logger:              mockLogger,
		runAsync:            func(f func()) { f() },
	}

	tests := []struct {
		name          string
		campaignID    string
		userEmail     string
		activityID    uint
		contributorID uint
		setupMocks    func()
		wantErr       bool
		expectedErr   string
		campaignKey   string
	}{
		{
			name:          "successful opt-out",
			campaignID:    "campaign1",
			userEmail:     "test@example.com",
			activityID:    1,
			campaignKey:   "campaign-key",
			contributorID: 1,
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{
								ID:         1,
								IsApproved: true,
								Contributors: []models.Contributor{
									{ID: 1, Email: "test@example.com"},
								},
							},
						},
						Contributors: []models.Contributor{
							{ID: 1, Email: "test@example.com"},
						},
					}, nil,
				)

				mockRepo.EXPECT().RemoveContributor(uint(1), uint(1)).Return(nil)

				mockBroadcaster.EXPECT().NewEvent(
					"campaign1",
					websocket.EventTypeActivityUpdated,
					mock.AnythingOfType("*models.Activity"),
				)
			},
			wantErr: false,
		},
		{
			name:          "opt-out fails - contributor not found",
			campaignID:    "campaign1",
			userEmail:     "test@example.com",
			activityID:    1,
			contributorID: 2,
			campaignKey:   "campaign-key",
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{ID: 1, IsApproved: true},
						},
						Contributors: []models.Contributor{
							{ID: 1, Email: "test@example.com"},
						},
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Contributor not found in this campaign",
		},
		{
			name:          "opt-out fails - wrong email",
			campaignID:    "campaign1",
			userEmail:     "wrong@example.com",
			activityID:    1,
			contributorID: 1,
			campaignKey:   "campaign-key",
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{ID: 1, IsApproved: true},
						},
						Contributors: []models.Contributor{
							{ID: 1, Email: "test@example.com"},
						},
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Only the contributor can perform this action",
		},
		{
			name:          "opt-out fails - already paid",
			campaignID:    "campaign1",
			userEmail:     "test@example.com",
			activityID:    1,
			campaignKey:   "campaign-key",
			contributorID: 1,
			setupMocks: func() {
				mockCampaign.EXPECT().GetCampaignByID("campaign1", "campaign-key").Return(
					&models.Campaign{
						ID: "campaign1",
						Activities: []models.Activity{
							{ID: 1, IsApproved: true},
						},
						Contributors: []models.Contributor{
							{
								ID:    1,
								Email: "test@example.com",
								Payment: &models.Payment{
									PaymentStatus: models.PaymentStatusSucceeded,
								},
							},
						},
					}, nil,
				)
			},
			wantErr:     true,
			expectedErr: "Action cannot be performed after making a payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCampaign.ExpectedCalls = nil
			tt.setupMocks()

			err := service.OptOutContributor(tt.campaignID, tt.userEmail, tt.campaignKey, tt.activityID, tt.contributorID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
