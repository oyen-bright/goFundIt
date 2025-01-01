package services

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	interfaces "github.com/oyen-bright/goFundIt/internal/services/mocks"
	logger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCronService(t *testing.T) {
	mockCampaignService := interfaces.NewMockCampaignService(t)
	mockNotificationService := interfaces.NewMockNotificationService(t)
	mockLogger := logger.NewMockLogger(t)

	cronService := NewCronService(mockCampaignService, mockNotificationService, mockLogger)

	t.Run("StartCronJobs", func(t *testing.T) {
		err := cronService.StartCronJobs()
		assert.NoError(t, err)
	})

	t.Run("StopCronJobs", func(t *testing.T) {
		cronService.StopCronJobs()
	})
}

func TestCleanUpExpiredCampaign(t *testing.T) {
	mockCampaignService := interfaces.NewMockCampaignService(t)
	mockNotificationService := interfaces.NewMockNotificationService(t)
	mockLogger := logger.NewMockLogger(t)

	expiredCampaign := models.Campaign{
		ID:      "test-id",
		EndDate: time.Now().Add(-24 * time.Hour), // Ended 24 hours ago
	}

	mockCampaignService.EXPECT().GetExpiredCampaigns().Return([]models.Campaign{expiredCampaign}, nil)
	mockCampaignService.EXPECT().GetCampaignByIDWithAllRelatedData(expiredCampaign.ID).Return(&expiredCampaign, nil)
	mockNotificationService.EXPECT().NotifyCampaignCleanUp(mock.Anything, mock.Anything).Return(nil)
	mockCampaignService.EXPECT().DeleteCampaign(expiredCampaign.ID).Return(nil)

	cronService := &cronService{
		campaignService:     mockCampaignService,
		notificationService: mockNotificationService,
		logger:              mockLogger,
	}
	cronService.cleanUpExpiredCampaign()

	// Allow some time for goroutines to complete
	time.Sleep(100 * time.Millisecond)
}

func TestCheckContributionReminders(t *testing.T) {
	mockCampaignService := interfaces.NewMockCampaignService(t)
	mockNotificationService := interfaces.NewMockNotificationService(t)
	mockLogger := logger.NewMockLogger(t)

	contributor := models.Contributor{
		ID:      1,
		Payment: nil, // Empty payment ID means not paid
	}

	activeCampaign := models.Campaign{
		ID:           "campaign-id",
		Contributors: []models.Contributor{contributor},
	}

	mockCampaignService.EXPECT().GetActiveCampaigns().Return([]models.Campaign{activeCampaign}, nil)
	mockNotificationService.EXPECT().SendContributionReminder(&contributor, &activeCampaign).Return(nil)

	cronService := &cronService{
		campaignService:     mockCampaignService,
		notificationService: mockNotificationService,
		logger:              mockLogger,
	}
	cronService.checkContributionReminders()

	time.Sleep(100 * time.Millisecond)
}

func TestCheckCampaignDeadline(t *testing.T) {
	mockCampaignService := interfaces.NewMockCampaignService(t)
	mockNotificationService := interfaces.NewMockNotificationService(t)
	mockLogger := logger.NewMockLogger(t)

	nearEndCampaign := models.Campaign{
		ID:      "campaign-id",
		EndDate: time.Now().Add(24 * time.Hour), // Ends in 24 hours
	}

	mockCampaignService.EXPECT().GetNearEndCampaigns().Return([]models.Campaign{nearEndCampaign}, nil)
	mockNotificationService.EXPECT().SendDeadlineReminder(&nearEndCampaign).Return(nil)

	cronService := &cronService{
		campaignService:     mockCampaignService,
		notificationService: mockNotificationService,
		logger:              mockLogger,
	}
	cronService.checkCampaignDeadline()

	//Allow time for go routine
	time.Sleep(100 * time.Millisecond)
}

func TestCreateJSONExport(t *testing.T) {
	campaign := models.Campaign{
		ID: "test-id",
	}

	filePath, err := createJSONExport(campaign)
	assert.NoError(t, err)
	assert.Contains(t, filePath, "campaign_export_test-id")
	assert.Contains(t, filePath, ".json")
}
