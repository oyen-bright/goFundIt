package services

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	mockInterfaces "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	serviceMocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	loggerMocks "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment(t *testing.T) {
	mockRepo := mockInterfaces.NewMockCommentRepository(t)
	mockAuth := serviceMocks.NewMockAuthService(t)
	mockActivity := serviceMocks.NewMockActivityService(t)
	mockNotification := serviceMocks.NewMockNotificationService(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewCommentService(mockRepo, mockAuth, mockActivity, mockNotification, mockBroadcaster, mockLogger)

	t.Run("successful comment creation", func(t *testing.T) {
		comment := &models.Comment{Content: "Test comment"}
		user := models.User{Handle: "testuser"}
		activity := models.Activity{ID: 1}

		mockAuth.On("GetUserByHandle", "testuser").Return(user, nil)
		mockActivity.On("GetActivityByID", uint(1), "campaign1").Return(activity, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.Comment")).Return(nil)
		mockBroadcaster.On("NewEvent", "campaign1", websocket.EventTypeCommentCreated, mock.AnythingOfType("*models.Comment")).Return(nil)
		mockNotification.On("NotifyCommentAddition", mock.AnythingOfType("*models.Comment"), &activity).Return(nil)

		err := service.CreateComment(comment, "campaign1", 1, "testuser")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
		mockActivity.AssertExpectations(t)
	})

	t.Run("invalid user", func(t *testing.T) {
		comment := &models.Comment{Content: "Test comment"}
		mockAuth.On("GetUserByHandle", "invalid").Return(models.User{}, errs.BadRequest("user not found", nil))

		err := service.CreateComment(comment, "campaign1", 1, "invalid")

		assert.Error(t, err)
		mockAuth.AssertExpectations(t)
	})
}

func TestCommentService_DeleteComment(t *testing.T) {
	mockRepo := mockInterfaces.NewMockCommentRepository(t)
	mockAuth := serviceMocks.NewMockAuthService(t)
	mockActivity := serviceMocks.NewMockActivityService(t)
	mockNotification := serviceMocks.NewMockNotificationService(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewCommentService(mockRepo, mockAuth, mockActivity, mockNotification, mockBroadcaster, mockLogger)

	t.Run("successful deletion", func(t *testing.T) {
		user := models.User{Handle: "testuser"}
		comment := models.Comment{ID: "123", CreatedByHandle: "testuser"}

		mockAuth.On("GetUserByHandle", "testuser").Return(user, nil)
		mockRepo.On("Get", "123").Return(comment, nil)
		mockRepo.On("Delete", "123").Return(nil)

		mockBroadcaster.On("NewEvent", "123", websocket.EventTypeCommentDeleted, "123").Return(nil)
		err := service.DeleteComment("123", "testuser")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("unauthorized deletion", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
		mockAuth.ExpectedCalls = nil
		mockAuth.Calls = nil
		mockActivity.ExpectedCalls = nil
		mockActivity.Calls = nil
		mockNotification.ExpectedCalls = nil
		mockNotification.Calls = nil
		mockBroadcaster.ExpectedCalls = nil
		mockBroadcaster.Calls = nil
		mockLogger.ExpectedCalls = nil
		mockLogger.Calls = nil

		user := models.User{Handle: "testuser"}
		comment := models.Comment{ID: "123", CreatedByHandle: "otheruser"}

		mockAuth.On("GetUserByHandle", "testuser").Return(user, nil)
		mockRepo.On("Get", "123").Return(comment, nil)

		err := service.DeleteComment("123", "testuser")

		mockRepo.AssertNotCalled(t, "Delete", "123")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})
}

func TestCommentService_GetActivityComments(t *testing.T) {
	mockRepo := mockInterfaces.NewMockCommentRepository(t)
	mockAuth := serviceMocks.NewMockAuthService(t)
	mockActivity := serviceMocks.NewMockActivityService(t)
	mockNotification := serviceMocks.NewMockNotificationService(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewCommentService(mockRepo, mockAuth, mockActivity, mockNotification, mockBroadcaster, mockLogger)

	t.Run("get comments successfully", func(t *testing.T) {
		expectedComments := []models.Comment{
			{ID: "1", Content: "Comment 1"},
			{ID: "2", Content: "Comment 2"},
		}

		mockRepo.On("GetByActivityID", uint(1)).Return(expectedComments, nil)

		comments, err := service.GetActivityComments(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedComments, comments)
		mockRepo.AssertExpectations(t)
	})
}

func TestCommentService_UpdateComment(t *testing.T) {
	mockRepo := mockInterfaces.NewMockCommentRepository(t)
	mockAuth := serviceMocks.NewMockAuthService(t)
	mockActivity := serviceMocks.NewMockActivityService(t)
	mockNotification := serviceMocks.NewMockNotificationService(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewCommentService(mockRepo, mockAuth, mockActivity, mockNotification, mockBroadcaster, mockLogger)

	t.Run("successful update", func(t *testing.T) {
		user := models.User{Handle: "testuser"}
		comment := models.Comment{
			ID:              "123",
			CreatedByHandle: "testuser",
			Content:         "Updated content",
		}

		mockAuth.On("GetUserByHandle", "testuser").Return(user, nil)
		mockRepo.On("Get", "123").Return(comment, nil)
		mockRepo.On("Update", &comment).Return(nil)

		err := service.UpdateComment(comment, "testuser")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("unauthorized update", func(t *testing.T) {

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
		mockAuth.ExpectedCalls = nil
		mockAuth.Calls = nil
		mockActivity.ExpectedCalls = nil
		mockActivity.Calls = nil
		mockNotification.ExpectedCalls = nil
		mockNotification.Calls = nil
		mockBroadcaster.ExpectedCalls = nil
		mockBroadcaster.Calls = nil
		mockLogger.ExpectedCalls = nil
		mockLogger.Calls = nil

		user := models.User{Handle: "testuser"}
		comment := models.Comment{
			ID:              "123",
			CreatedByHandle: "otheruser",
			Content:         "Updated content",
		}

		mockAuth.On("GetUserByHandle", "testuser").Return(user, nil)
		mockRepo.On("Get", "123").Return(comment, nil)

		err := service.UpdateComment(comment, "testuser")

		// Assert that an error is returned.
		assert.Error(t, err)

		// Assert that Update was not called.
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)

		// Verify expectations on the other mocks.
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

}

func TestCommentService_GetCommentReplies(t *testing.T) {
	mockRepo := mockInterfaces.NewMockCommentRepository(t)
	mockAuth := serviceMocks.NewMockAuthService(t)
	mockActivity := serviceMocks.NewMockActivityService(t)
	mockNotification := serviceMocks.NewMockNotificationService(t)
	mockBroadcaster := serviceMocks.NewMockEventBroadcaster(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewCommentService(mockRepo, mockAuth, mockActivity, mockNotification, mockBroadcaster, mockLogger)

	t.Run("get replies successfully", func(t *testing.T) {
		expectedReplies := []models.Comment{
			{ID: "2", Content: "Reply 1", ParentID: stringPtr("1")},
			{ID: "3", Content: "Reply 2", ParentID: stringPtr("1")},
		}

		mockRepo.On("FindReplies", "1").Return(expectedReplies, nil)

		replies, err := service.GetCommentReplies("1")

		assert.NoError(t, err)
		assert.Equal(t, expectedReplies, replies)
		mockRepo.AssertExpectations(t)
	})
}

func stringPtr(s string) *string {
	return &s
}
