package services

import (
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
	mockAuth "github.com/oyen-bright/goFundIt/internal/services/mocks"
	mockEmailer "github.com/oyen-bright/goFundIt/pkg/email/mocks"
	mockFCM "github.com/oyen-bright/goFundIt/pkg/fcm/mocks"

	mockLogger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (*notificationService, *mockEmailer.MockEmailer, *mockFCM.MockFCM, *mockAuth.MockAuthService) {
	mockEmailer := mockEmailer.NewMockEmailer(t)
	mockFCMClient := mockFCM.NewMockFCM(t)
	mockAuthService := mockAuth.NewMockAuthService(t)
	mockLogger := mockLogger.NewMockLogger(t)

	service := NewNotificationService(
		mockEmailer,
		mockAuthService,
		mockFCMClient,
		mockLogger,
	)

	return service.(*notificationService), mockEmailer, mockFCMClient, mockAuthService
}

func TestNotifyActivityAddition(t *testing.T) {
	service, mockEmailer, _, _ := setupTest(t)

	activity := &models.Activity{
		Title:    "Test Activity",
		Subtitle: "Test Subtitle",
		Cost:     100,
	}

	campaign := &models.Campaign{
		ID: "campaign123",
		Contributors: []models.Contributor{
			{Email: "test1@example.com"},
			{Email: "test2@example.com"},
		},
	}

	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil)

	err := service.NotifyActivityAddition(activity, campaign)

	assert.NoError(t, err)
	mockEmailer.AssertExpectations(t)
}

func TestNotifyActivityApproved(t *testing.T) {
	service, mockEmailer, _, _ := setupTest(t)

	activity := &models.Activity{
		Title:     "Test Activity",
		Subtitle:  "Test Subtitle",
		UpdatedAt: time.Now(),
	}

	campaign := &models.Campaign{
		ID: "campaign123",
		CreatedBy: models.User{
			Email: "creator@example.com",
		},
		Contributors: []models.Contributor{
			{Email: "test1@example.com"},
		},
	}

	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil)

	err := service.NotifyActivityApproved(activity, campaign)

	assert.NoError(t, err)
	mockEmailer.AssertExpectations(t)
}

func TestNotifyPaymentReceived(t *testing.T) {
	service, mockEmailer, mockFCM, _ := setupTest(t)

	fcmToken := "test-fcm-token"
	contributor := &models.Contributor{
		Email: "contributor@example.com",
		Name:  "Test Contributor",
		Payment: &models.Payment{
			Amount: 100,
		},
	}

	campaign := &models.Campaign{
		ID: "campaign123",
		CreatedBy: models.User{
			FCMToken: &fcmToken,
		},
	}

	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil)
	mockFCM.On("SendNotification", mock.Anything, fcmToken, mock.AnythingOfType("fcm.NotificationData")).Return(nil)

	err := service.NotifyPaymentReceived(contributor, campaign)

	assert.NoError(t, err)
	mockEmailer.AssertExpectations(t)
	mockFCM.AssertExpectations(t)
}

func TestNotifyCampaignCreation(t *testing.T) {
	service, mockEmailer, _, _ := setupTest(t)

	campaign := &models.Campaign{
		ID:          "campaign123",
		Title:       "Test Campaign",
		Description: "Test Description",
		Key:         "test-key",
		CreatedBy: models.User{
			Email: "creator@example.com",
		},
		Contributors: []models.Contributor{
			{
				Email: "contributor1@example.com",
				Name:  "Contributor 1",
			},
		},
		Activities: []models.Activity{
			{
				Title:    "Activity 1",
				Subtitle: "Subtitle 1",
			},
		},
	}

	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil).Times(2)

	err := service.NotifyCampaignCreation(campaign)

	assert.NoError(t, err)
	// Wait for goroutine to complete
	time.Sleep(100 * time.Millisecond)
	mockEmailer.AssertExpectations(t)
}

func TestSendSystemNotification(t *testing.T) {
	service, mockEmailer, _, mockAuth := setupTest(t)

	users := []models.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}

	mockAuth.On("GetAllUser").Return(users, nil)
	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil)

	err := service.SendSystemNotification("TEST", "Test message")

	assert.NoError(t, err)
	mockEmailer.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestNotifyCommentAddition(t *testing.T) {
	service, mockEmailer, _, _ := setupTest(t)

	comment := &models.Comment{
		Content: "Test comment",
		CreatedBy: models.User{
			Handle: "testuser",
		},
	}

	activity := &models.Activity{
		Title:      "Test Activity",
		CampaignID: "campaign123",
		CreatedBy: models.User{
			Email: "creator@example.com",
		},
		Contributors: []models.Contributor{
			{Email: "contributor@example.com"},
		},
	}

	mockEmailer.On("SendEmailTemplate", mock.AnythingOfType("email.EmailTemplate")).Return(nil)

	err := service.NotifyCommentAddition(comment, activity)

	assert.NoError(t, err)
	mockEmailer.AssertExpectations(t)
}
