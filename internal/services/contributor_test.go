package services

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	mockRepo "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	mockService "github.com/oyen-bright/goFundIt/internal/services/mocks"
	loggerMock "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupContributorTest(t *testing.T) (
	*mockRepo.MockContributorRepository,
	*mockService.MockCampaignService,
	*mockService.MockAnalyticsService,
	*mockService.MockAuthService,
	*mockService.MockNotificationService,
	*mockService.MockEventBroadcaster,
	*loggerMock.MockLogger,
	*contributorService,
) {
	repo := mockRepo.NewMockContributorRepository(t)
	campaignService := mockService.NewMockCampaignService(t)
	analyticsService := mockService.NewMockAnalyticsService(t)
	authService := mockService.NewMockAuthService(t)
	notificationService := mockService.NewMockNotificationService(t)
	broadcaster := mockService.NewMockEventBroadcaster(t)
	logger := loggerMock.NewMockLogger(t)

	service := &contributorService{
		repo:                repo,
		campaignService:     campaignService,
		analyticsService:    analyticsService,
		authService:         authService,
		notificationService: notificationService,
		broadcaster:         broadcaster,
		logger:              logger,
		runAsync:            func(f func()) { f() },
	}

	return repo, campaignService, analyticsService, authService, notificationService, broadcaster, logger, service
}

func TestAddContributorToCampaign(t *testing.T) {
	repo, campaignService, _, authService, notificationService, broadcaster, _, service := setupContributorTest(t)

	testCases := []struct {
		name        string
		contributor *models.Contributor
		campaignID  string
		campaignKey string
		userHandle  string

		setupMocks    func()
		expectedError bool
	}{
		{
			name: "Success - New User",
			contributor: &models.Contributor{
				Name:  "Test User",
				Email: "test@example.com",
			},
			campaignID:  "campaign-123",
			campaignKey: "key-123",
			userHandle:  "creator",
			setupMocks: func() {
				campaign := &models.Campaign{
					ID:      "campaign-123",
					EndDate: time.Now().AddDate(0, 0, 30),
					CreatedBy: models.User{
						Handle: "creator",
					},
				}

				campaignService.EXPECT().GetCampaignByID("campaign-123", "key-123").Return(campaign, nil)
				authService.EXPECT().FindUserByEmail("test@example.com").Return(nil, nil)
				authService.EXPECT().CreateUser(mock.AnythingOfType("models.User")).Return(nil)
				repo.EXPECT().Create(mock.AnythingOfType("*models.Contributor")).Return(nil)
				broadcaster.EXPECT().NewEvent("campaign-123", websocket.EventTypeContributionCreated, mock.Anything)
				notificationService.EXPECT().NotifyContributorAdded(mock.Anything, mock.Anything).Return(nil)
				campaignService.EXPECT().RecalculateTargetAmount("campaign-123")
			},
			expectedError: false,
		},
		{
			name:        "Failure - Campaign Not Found",
			campaignKey: "key-123",
			contributor: &models.Contributor{
				Name:  "Test User",
				Email: "test@example.com",
			},
			campaignID: "invalid-id",
			setupMocks: func() {
				campaignService.EXPECT().GetCampaignByID("invalid-id", "key-123").Return(nil, assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		campaignService.ExpectedCalls = nil
		campaignService.Calls = nil
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := service.AddContributorToCampaign(tc.contributor, tc.campaignID, tc.campaignKey, tc.userHandle)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateContributor(t *testing.T) {
	repo, campaignService, _, _, _, broadcaster, mockLogger, service := setupContributorTest(t)

	testCases := []struct {
		name          string
		contributor   *models.Contributor
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "Success",
			contributor: &models.Contributor{
				CampaignID: "campaign-123",
				Name:       "Updated Name",
			},
			setupMocks: func() {
				repo.EXPECT().Update(mock.AnythingOfType("*models.Contributor")).Return(nil)
				broadcaster.EXPECT().NewEvent("campaign-123", websocket.EventTypeContributorUpdated, mock.Anything)

				campaignService.EXPECT().RecalculateTargetAmount("campaign-123")

			},
			expectedError: false,
		},
		{
			name: "Failure - Update Error",
			contributor: &models.Contributor{
				CampaignID: "campaign-123",
			},
			setupMocks: func() {
				repo.EXPECT().Update(mock.AnythingOfType("*models.Contributor")).Return(assert.AnError)
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		repo.ExpectedCalls = nil
		repo.Calls = nil
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := service.UpdateContributor(tc.contributor)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetContributorByID(t *testing.T) {
	repo, _, _, _, _, _, mockLogger, service := setupContributorTest(t)

	testCases := []struct {
		name          string
		contributorID uint
		setupMocks    func()
		expectedError bool
	}{
		{
			name:          "Success",
			contributorID: 1,
			setupMocks: func() {
				repo.EXPECT().GetContributorById(uint(1), true).Return(models.Contributor{ID: 1}, nil)
			},
			expectedError: false,
		},
		{
			name:          "Failure - Not Found",
			contributorID: 999,
			setupMocks: func() {
				repo.EXPECT().GetContributorById(uint(999), true).Return(models.Contributor{}, assert.AnError)
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			contributor, err := service.GetContributorByID(tc.contributorID)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.contributorID, contributor.ID)
			}
		})
	}
}

func TestGetContributorsByCampaignID(t *testing.T) {
	repo, _, _, _, _, _, mockLogger, service := setupContributorTest(t)

	testCases := []struct {
		name          string
		campaignID    string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:       "Success",
			campaignID: "campaign-123",
			setupMocks: func() {
				contributors := []models.Contributor{{ID: 1}, {ID: 2}}
				repo.EXPECT().GetContributorsByCampaignID("campaign-123").Return(contributors, nil)
			},
			expectedError: false,
		},
		{
			name:       "Failure",
			campaignID: "invalid-id",
			setupMocks: func() {
				repo.EXPECT().GetContributorsByCampaignID("invalid-id").Return(nil, assert.AnError)
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			contributors, err := service.GetContributorsByCampaignID(tc.campaignID)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, contributors)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, contributors)
				assert.Len(t, contributors, 2)
			}
		})
	}
}

func TestRemoveContributorFromCampaign(t *testing.T) {
	repo, campaignService, _, _, _, broadcaster, _, service := setupContributorTest(t)

	testCases := []struct {
		name          string
		contributorID uint
		campaignID    string
		userHandle    string
		campaignKey   string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:          "Success",
			contributorID: 1,
			campaignID:    "campaign-123",
			userHandle:    "creator",
			campaignKey:   "key-123",
			setupMocks: func() {
				contributor := &models.Contributor{ID: 1, CampaignID: "campaign-123"}
				campaign := &models.Campaign{
					ID: "campaign-123",
					CreatedBy: models.User{
						Handle: "creator",
					},
					Contributors: []models.Contributor{*contributor},
				}

				campaignService.EXPECT().GetCampaignByID("campaign-123", "key-123").Return(campaign, nil)
				repo.EXPECT().Delete(contributor).Return(nil)
				broadcaster.EXPECT().NewEvent("campaign-123", websocket.EventTypeContributorDeleted, mock.Anything)
				campaignService.EXPECT().RecalculateTargetAmount("campaign-123")
			},
			expectedError: false,
		},
		{
			name:          "Failure - Unauthorized",
			contributorID: 1,
			campaignID:    "campaign-123",
			userHandle:    "not-creator",
			campaignKey:   "key-123",
			setupMocks: func() {
				campaign := &models.Campaign{
					ID: "campaign-123",
					CreatedBy: models.User{
						Handle: "creator",
					},
				}
				campaignService.EXPECT().GetCampaignByID("campaign-123", "key-123").Return(campaign, nil)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := service.RemoveContributorFromCampaign(tc.contributorID, tc.campaignID, tc.userHandle, tc.campaignKey)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
