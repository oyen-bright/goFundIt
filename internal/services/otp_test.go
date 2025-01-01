package services

import (
	"strings"
	"testing"
	"time"

	"github.com/oyen-bright/goFundIt/pkg/email"

	"github.com/oyen-bright/goFundIt/internal/models"
	repoMocks "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	emailMocks "github.com/oyen-bright/goFundIt/pkg/email/mocks"
	encryptorMocks "github.com/oyen-bright/goFundIt/pkg/encryption/mocks"
	loggerMocks "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestOTP(t *testing.T) {
	// Setup mocks
	mockRepo := repoMocks.NewMockOTPRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	// Create service with mocks
	service := &otpService{
		repo:     mockRepo,
		emailer:  mockEmailer,
		logger:   mockLogger,
		runAsync: false,
	}

	tests := []struct {
		name          string
		email         string
		userName      string
		setupMocks    func()
		expectedError bool
	}{
		{
			name:     "Success case",
			email:    "test@example.com",
			userName: "Test User",
			setupMocks: func() {

				mockRepo.EXPECT().Add(mock.Anything).Return(nil)
				mockRepo.EXPECT().InvalidateOtherOTPs(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mockEmailer.EXPECT().SendEmailTemplate(
					mock.MatchedBy(func(et email.EmailTemplate) bool {
						emailMatch := len(et.To) == 1 && et.To[0] == "test@example.com"
						subjectMatch := et.Subject == "Email Verification - GoFund It"
						pathMatch := strings.HasSuffix(et.Path, "/templates/personal/email_verification.html")
						dataMatch := et.Data != nil &&
							et.Data["name"] == "Test User" &&
							et.Data["verificationCode"] != nil

						return emailMatch && subjectMatch && pathMatch && dataMatch
					}),
				).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks for each test case
			tt.setupMocks()

			otp, err := service.RequestOTP(tt.email, tt.userName)

			// Verify results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, otp.Code)
				assert.Equal(t, tt.email, otp.Email)
				assert.Equal(t, tt.userName, otp.Name)
			}
		})
	}
}

func TestVerifyOTP(t *testing.T) {
	// Setup mocks
	mockRepo := repoMocks.NewMockOTPRepository(t)
	mockEmailer := emailMocks.NewMockEmailer(t)
	mockEncryptor := encryptorMocks.NewMockEncryptor(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	service := NewOTPService(mockRepo, mockEmailer, mockLogger)

	validOTP := &models.Otp{
		Email:     "test@example.com",
		Code:      "123456",
		RequestId: "test-request-id",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	expiredOTP := &models.Otp{
		Email:     "test1@example.com",
		Code:      "123456",
		RequestId: "test-request-id",
		ExpiresAt: time.Now().Add(-time.Hour),
		Name:      "Test User",
	}

	tests := []struct {
		name          string
		email         string
		code          string
		requestID     string
		mockSetup     func()
		expectedError bool
		errorMessage  string
	}{
		{
			name:      "Valid OTP",
			email:     "test@example.com",
			code:      "123456",
			requestID: "test-request-id",
			mockSetup: func() {

				// Mock repository lookup
				mockRepo.EXPECT().GetByEmailAndReference("test@example.com", "test-request-id").Return(validOTP, nil)

				// Mock deletion after successful verification
				mockRepo.EXPECT().Delete(validOTP).Return(nil)
			},
			expectedError: false,
		},
		{
			name:      "Expired OTP",
			email:     "test1@example.com",
			code:      "123456",
			requestID: "test-request-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetByEmailAndReference("test1@example.com", "test-request-id").Return(expiredOTP, nil)
			},
			expectedError: true,
			errorMessage:  "OTP has expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockEncryptor.ExpectedCalls = nil
			mockEmailer.ExpectedCalls = nil

			tt.mockSetup()

			_, err := service.VerifyOTP(tt.email, tt.code, tt.requestID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.code, validOTP.Code)
				assert.Equal(t, tt.requestID, validOTP.RequestId)
			}
		})
	}
}
