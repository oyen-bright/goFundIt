package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	dto "github.com/oyen-bright/goFundIt/internal/api/dto/campaign"
	"github.com/oyen-bright/goFundIt/internal/models"
	mockRepo "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	mockInterfaces "github.com/oyen-bright/goFundIt/internal/services/mocks"
	encrypt "github.com/oyen-bright/goFundIt/pkg/encryption/mocks"
	mockLogger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCampaignService(t *testing.T) (
	*campaignService,
	*mockRepo.MockCampaignRepository,
	*mockInterfaces.MockAuthService,
	*mockInterfaces.MockAnalyticsService,
	*mockInterfaces.MockNotificationService,
	*mockInterfaces.MockEventBroadcaster,
	*mockLogger.MockLogger,
	*encrypt.MockEncryptor,
) {
	mockRepo := mockRepo.NewMockCampaignRepository(t)
	mockAuth := mockInterfaces.NewMockAuthService(t)
	mockAnalytics := mockInterfaces.NewMockAnalyticsService(t)
	mockNotification := mockInterfaces.NewMockNotificationService(t)
	mockBroadcaster := mockInterfaces.NewMockEventBroadcaster(t)
	mockLogger := mockLogger.NewMockLogger(t)
	mockEncryptor := encrypt.NewMockEncryptor(t)

	service := &campaignService{
		encryptor:           mockEncryptor,
		repo:                mockRepo,
		authService:         mockAuth,
		analyticsService:    mockAnalytics,
		notificationService: mockNotification,
		broadcaster:         mockBroadcaster,
		logger:              mockLogger,
		runAsync:            func(f func()) { f() },
	}

	return service, mockRepo, mockAuth, mockAnalytics, mockNotification, mockBroadcaster, mockLogger, mockEncryptor
}

func TestCreateCampaign(t *testing.T) {
	service, mockRepo, mockAuth, mockAnalytics, mockNotification, _, mockLogger, encryptor := setupCampaignService(t)

	t.Run("successful campaign creation", func(t *testing.T) {
		campaignKey := "test_key"
		// Test data
		userHandle := "test_user"
		campaign := &models.Campaign{
			Key:         campaignKey,
			Title:       "Test Campaign",
			Description: "Test Description",
			Contributors: []models.Contributor{
				{Email: "test@example.com"},
			},
		}
		user := models.User{
			Handle: userHandle,
			Email:  "test@example.com",
		}

		// Setup expectations
		mockRepo.EXPECT().GetByHandle(userHandle).Return(models.Campaign{}, nil)
		mockAuth.EXPECT().FindExistingAndNonExistingUsers([]string{"test@example.com"}).
			Return([]models.User{}, []string{"test@example.com"}, nil)
		mockAuth.EXPECT().CreateUsers(mock.AnythingOfType("[]models.User")).Return([]models.User{}, nil)
		mockAuth.EXPECT().GetUserByHandle(userHandle).Return(user, nil)

		encryptor.EXPECT().EncryptStruct(mock.AnythingOfType("*models.Campaign"), mock.AnythingOfType("string")).Return(mock.AnythingOfType("*models.Campaign"), nil)
		mockRepo.EXPECT().Create(campaign).Return(*campaign, nil)

		encryptor.EXPECT().DecryptStruct(mock.AnythingOfType("*models.Campaign"), mock.AnythingOfType("string")).Return(mock.Anything, nil)

		platformAnalytics := &models.PlatformAnalytics{}
		mockAnalytics.EXPECT().GetCurrentData().Return(platformAnalytics)
		mockNotification.EXPECT().NotifyCampaignCreation(campaign).Return(nil)

		// Execute
		result, err := service.CreateCampaign(campaign, userHandle)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, campaign.Title, result.Title)
	})

	mockRepo.ExpectedCalls = nil
	mockNotification.ExpectedCalls = nil

	t.Run("error - user already has campaign", func(t *testing.T) {

		userHandle := "test_user"
		campaign := &models.Campaign{Title: "Test Campaign"}
		existingCampaign := models.Campaign{ID: "existing-id"}

		mockRepo.EXPECT().GetByHandle(userHandle).Return(existingCampaign, errors.New("You already have an active campaign"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		_, err := service.CreateCampaign(campaign, userHandle)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "have an active campaign")
	})

	t.Run("error - invalid contributors", func(t *testing.T) {
		userHandle := "test_user"
		campaign := &models.Campaign{
			Contributors: []models.Contributor{{Email: "test@example.com"}},
		}
		existingUser := models.User{Email: "test@example.com", Contributions: []models.Contributor{{CampaignID: "other-campaign"}}}

		mockRepo.EXPECT().GetByHandle(userHandle).Return(models.Campaign{}, nil)
		mockAuth.EXPECT().FindExistingAndNonExistingUsers([]string{"test@example.com"}).
			Return([]models.User{existingUser}, []string{}, nil)

		_, err := service.CreateCampaign(campaign, userHandle)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "have an active campaign")
	})
}

func TestUpdateCampaign(t *testing.T) {
	service, mockRepo, _, _, mockNotification, mockBroadcaster, mockLogger, encryptor := setupCampaignService(t)

	t.Run("successful campaign update", func(t *testing.T) {
		campaignID := "test_id"
		campaignKey := "test_key"
		userHandle := "test_user"

		updatedTitle := "Updated Title"
		updatedDescription := "Updated Description"
		endDate := time.Now().Add(24 * time.Hour)

		updateReq := dto.CampaignUpdateRequest{
			Title:       &updatedTitle,
			Description: &updatedDescription,
			EndDate:     &endDate,
		}

		existingCampaign := &models.Campaign{
			ID:  campaignID,
			Key: campaignKey,
			CreatedBy: models.User{
				Handle: userHandle,
			},
		}

		// Setup expectations
		mockRepo.EXPECT().GetByID(campaignID).Return(*existingCampaign, nil)
		encryptor.EXPECT().EncryptStruct(mock.AnythingOfType("*models.Campaign"), campaignKey).Return(mock.AnythingOfType("*models.Campaign"), nil)

		mockRepo.EXPECT().Update(mock.AnythingOfType("*models.Campaign")).Return(*existingCampaign, nil)
		encryptor.EXPECT().DecryptStruct(mock.AnythingOfType("*models.Campaign"), campaignKey).Return(mock.Anything, nil)

		mockBroadcaster.EXPECT().NewEvent(campaignID, mock.Anything, mock.Anything)
		mockNotification.EXPECT().NotifyCampaignUpdate(mock.AnythingOfType("*models.Campaign"), "").Return(nil)

		// Execute
		result, err := service.UpdateCampaign(updateReq, campaignID, campaignKey, userHandle)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("error - campaign not found", func(t *testing.T) {
		updatedTitle := "Updated Title"

		updateReq := dto.CampaignUpdateRequest{Title: &updatedTitle}
		mockRepo.EXPECT().GetByID("non-existent").Return(models.Campaign{}, fmt.Errorf("not found"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		_, err := service.UpdateCampaign(updateReq, "non-existent", "campaign-key", "user123")
		assert.Error(t, err)
	})

	t.Run("error - unauthorized update", func(t *testing.T) {
		campaignID := "test-id"
		updatedTitle := "Updated Title"
		campaignKey := "test_key"

		updateReq := dto.CampaignUpdateRequest{Title: &updatedTitle}
		existingCampaign := models.Campaign{
			ID:        campaignID,
			Key:       campaignKey,
			CreatedBy: models.User{Handle: "different-user"},
		}

		mockRepo.EXPECT().GetByID(campaignID).Return(existingCampaign, nil)
		encryptor.EXPECT().DecryptStruct(mock.AnythingOfType("*models.Campaign"), campaignKey).Return(mock.Anything, nil)

		_, err := service.UpdateCampaign(updateReq, campaignID, campaignKey, "unauthorized-user")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unauthorized")
	})
}

func TestGetCampaignByID(t *testing.T) {
	service, mockRepo, _, _, _, _, _, encryptor := setupCampaignService(t)

	t.Run("successful campaign retrieval", func(t *testing.T) {
		campaignID := "test_id"
		campaignKey := "test_key"
		expectedCampaign := models.Campaign{
			ID:    campaignID,
			Title: "Test Campaign",
		}

		mockRepo.EXPECT().GetByID(campaignID).Return(expectedCampaign, nil)
		encryptor.EXPECT().DecryptStruct(mock.AnythingOfType("*models.Campaign"), campaignKey).Return(mock.Anything, nil)

		result, err := service.GetCampaignByID(campaignID, campaignKey)

		assert.NoError(t, err)
		assert.Equal(t, expectedCampaign.ID, result.ID)
		assert.Equal(t, expectedCampaign.Title, result.Title)
	})

}

func TestGetActiveCampaigns(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := setupCampaignService(t)

	t.Run("successful active campaigns retrieval", func(t *testing.T) {
		expectedCampaigns := []models.Campaign{
			{ID: "1", Title: "Campaign 1"},
			{ID: "2", Title: "Campaign 2"},
		}

		mockRepo.EXPECT().GetActiveCampaigns().Return(expectedCampaigns, nil)

		results, err := service.GetActiveCampaigns()

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, expectedCampaigns[0].ID, results[0].ID)
		assert.Equal(t, expectedCampaigns[1].ID, results[1].ID)
	})

}

func TestDeleteCampaign(t *testing.T) {
	service, mockRepo, _, _, _, _, mockLogger, _ := setupCampaignService(t)

	t.Run("successful deletion", func(t *testing.T) {
		campaignID := "test-id"
		mockRepo.EXPECT().Delete(campaignID).Return(nil)

		err := service.DeleteCampaign(campaignID)
		assert.NoError(t, err)
	})

	t.Run("error - deletion failed", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil

		campaignID := "test-id"
		mockRepo.EXPECT().Delete(campaignID).Return(fmt.Errorf("deletion error"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		err := service.DeleteCampaign(campaignID)
		assert.Error(t, err)
	})
}

func TestGetCampaignByIDWithContributors(t *testing.T) {
	service, mockRepo, _, _, _, _, mockLogger, _ := setupCampaignService(t)

	t.Run("successful retrieval with contributors", func(t *testing.T) {
		campaignID := "test-id"
		expectedCampaign := models.Campaign{
			ID: campaignID,
			Contributors: []models.Contributor{
				{Email: "contributor@example.com"},
			},
		}

		mockRepo.EXPECT().GetByIDWithSelectedData(campaignID,
			models.PreloadOption{Contributors: true, Payout: true}).
			Return(expectedCampaign, nil)

		result, err := service.GetCampaignByIDWithContributors(campaignID)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedCampaign.Contributors), len(result.Contributors))
	})

	t.Run("error - campaign not found", func(t *testing.T) {
		campaignID := "non-existent"
		mockRepo.EXPECT().GetByIDWithSelectedData(campaignID, mock.Anything).
			Return(models.Campaign{}, fmt.Errorf("not found"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		_, err := service.GetCampaignByIDWithContributors(campaignID)
		assert.Error(t, err)
	})
}

func TestGetCampaignByIDWithAllRelatedData(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := setupCampaignService(t)

	t.Run("successful retrieval with all data", func(t *testing.T) {
		campaignID := "test-id"
		expectedCampaign := models.Campaign{
			ID:           campaignID,
			Contributors: []models.Contributor{{Email: "test@example.com"}},
			Activities:   []models.Activity{{Title: "Test Activity"}},
		}

		mockRepo.EXPECT().GetByIDWithSelectedData(campaignID, mock.Anything).
			Return(expectedCampaign, nil)

		result, err := service.GetCampaignByIDWithAllRelatedData(campaignID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCampaign.ID, result.ID)
		assert.NotEmpty(t, result.Contributors)
		assert.NotEmpty(t, result.Activities)
	})
}

func TestGetExpiredCampaigns(t *testing.T) {
	service, mockRepo, _, _, _, _, mockLogger, _ := setupCampaignService(t)

	t.Run("successful retrieval of expired campaigns", func(t *testing.T) {
		expectedCampaigns := []models.Campaign{
			{ID: "1", EndDate: time.Now().Add(-24 * time.Hour)},
			{ID: "2", EndDate: time.Now().Add(-48 * time.Hour)},
		}

		mockRepo.EXPECT().GetExpiredCampaigns().Return(expectedCampaigns, nil)

		results, err := service.GetExpiredCampaigns()
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("error - retrieval failed", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil

		mockRepo.EXPECT().GetExpiredCampaigns().Return(nil, fmt.Errorf("database error"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		_, err := service.GetExpiredCampaigns()
		assert.Error(t, err)
	})
}

func TestGetNearEndCampaigns(t *testing.T) {
	service, mockRepo, _, _, _, _, mockLogger, _ := setupCampaignService(t)

	t.Run("successful retrieval of near-end campaigns", func(t *testing.T) {
		expectedCampaigns := []models.Campaign{
			{ID: "1", EndDate: time.Now().Add(24 * time.Hour)},
			{ID: "2", EndDate: time.Now().Add(48 * time.Hour)},
		}

		mockRepo.EXPECT().GetNearEndCampaigns().Return(expectedCampaigns, nil)

		results, err := service.GetNearEndCampaigns()
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("error - retrieval failed", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.EXPECT().GetNearEndCampaigns().Return(nil, fmt.Errorf("database error"))
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

		_, err := service.GetNearEndCampaigns()
		assert.Error(t, err)
	})
}

func TestRecalculateTargetAmount(t *testing.T) {
	service, mockRepo, _, _, _, mockBroadcaster, _, _ := setupCampaignService(t)

	t.Run("successful recalculation", func(t *testing.T) {
		campaignID := "test-id"
		campaign := models.Campaign{
			ID: campaignID,
			Contributors: []models.Contributor{
				{Amount: 100},
				{Amount: 200},
			},
		}

		mockRepo.EXPECT().GetByID(campaignID).Return(campaign, nil)
		mockRepo.EXPECT().Update(mock.AnythingOfType("*models.Campaign")).Return(campaign, nil)
		mockBroadcaster.EXPECT().NewEvent(campaignID, mock.Anything, mock.Anything)

		service.RecalculateTargetAmount(campaignID)
		mockRepo.AssertExpectations(t)
		mockBroadcaster.AssertExpectations(t)
	})

	t.Run("error - campaign not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		campaignID := "non-existent"
		mockRepo.EXPECT().GetByID(campaignID).Return(models.Campaign{}, fmt.Errorf("not found"))

		service.RecalculateTargetAmount(campaignID)
	})
}
