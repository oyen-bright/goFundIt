package services

import (
	"errors"
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	mocks "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	emailMocks "github.com/oyen-bright/goFundIt/pkg/email/mocks"
	loggerMocks "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAnalyticsService(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	// Mock getCurrentData behavior
	mockRepo.EXPECT().Get(mock.Anything).Return(&models.PlatformAnalytics{}, nil)

	// Create service
	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	// Assert
	assert.NotNil(t, service, "Service should not be nil")
}

func TestStartAnalytics(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	mockRepo.EXPECT().Get(mock.Anything).Return(&models.PlatformAnalytics{}, nil)
	mockLogger.EXPECT().Info(mock.Anything, mock.Anything).Return()

	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	err := service.StartAnalytics()

	assert.NoError(t, err, "StartAnalytics should not return an error")
}

func TestProcessAnalyticsNow(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	currentData := &models.PlatformAnalytics{
		TotalUsers:     100,
		TotalCampaigns: 50,
	}

	yesterdayData := &models.PlatformAnalytics{
		TotalUsers:     90,
		TotalCampaigns: 45,
	}

	mockRepo.EXPECT().Get(mock.Anything).Return(currentData, nil)
	mockRepo.EXPECT().Save(mock.Anything).Return(nil)
	mockRepo.EXPECT().Get(mock.Anything).Return(yesterdayData, nil)
	mockEmailer.EXPECT().SendEmailTemplate(mock.Anything).Return(nil)
	mockLogger.EXPECT().Info(mock.Anything, mock.Anything).Return()

	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	err := service.ProcessAnalyticsNow()

	assert.NoError(t, err, "ProcessAnalyticsNow should not return an error")
}

func TestProcessAnalyticsNow_Error(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	// Set up error case
	mockRepo.EXPECT().Get(mock.Anything).Return(&models.PlatformAnalytics{}, nil)
	mockRepo.EXPECT().Save(mock.Anything).Return(errors.New("database error"))
	mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	// Test
	err := service.ProcessAnalyticsNow()

	// Assert
	assert.Error(t, err, "ProcessAnalyticsNow should return an error")
}

func TestGetCurrentData(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	expectedData := &models.PlatformAnalytics{
		TotalUsers:     100,
		TotalCampaigns: 50,
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().Get(mock.Anything).Return(expectedData, nil)

	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	data := service.GetCurrentData()

	assert.Equal(t, expectedData, data, "GetCurrentData should return the expected data")
}

func TestStopAnalytics(t *testing.T) {
	mockRepo := mocks.NewMockAnalyticsRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)
	testEmail := "admin@test.com"

	mockRepo.EXPECT().Get(mock.Anything).Return(&models.PlatformAnalytics{}, nil)
	mockLogger.EXPECT().Info(mock.Anything, mock.Anything).Return()

	service := NewAnalyticsService(mockRepo, testEmail, mockEmailer, mockLogger)

	_ = service.StartAnalytics()

	service.StopAnalytics()

	assert.True(t, true, "StopAnalytics should not panic")
}
