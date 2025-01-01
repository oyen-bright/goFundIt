package services

import (
	"errors"
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	mock_interfaces "github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	mock_services "github.com/oyen-bright/goFundIt/internal/services/mocks"
	mock_encryption "github.com/oyen-bright/goFundIt/pkg/encryption/mocks"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	mock_jwt "github.com/oyen-bright/goFundIt/pkg/jwt/mocks"
	mock_logger "github.com/oyen-bright/goFundIt/pkg/logger/mocks"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthTest(t *testing.T) (
	*mock_interfaces.MockAuthRepository,
	*mock_services.MockOTPService,
	*mock_services.MockAnalyticsService,
	*mock_encryption.MockEncryptor,
	*mock_jwt.MockJwt,
	*mock_logger.MockLogger,
	*authService,
) {
	mockAuthRepo := mock_interfaces.NewMockAuthRepository(t)
	mockOtpService := mock_services.NewMockOTPService(t)
	mockAnalyticsService := mock_services.NewMockAnalyticsService(t)
	mockEncryptor := mock_encryption.NewMockEncryptor(t)
	mockJwt := mock_jwt.NewMockJwt(t)
	mockLogger := mock_logger.NewMockLogger(t)

	service := &authService{
		authRepo:         mockAuthRepo,
		otpService:       mockOtpService,
		analyticsService: mockAnalyticsService,
		encryptor:        mockEncryptor,
		jwt:              mockJwt,
		logger:           mockLogger,
	}

	return mockAuthRepo, mockOtpService, mockAnalyticsService, mockEncryptor, mockJwt, mockLogger, service
}

func TestRequestAuth(t *testing.T) {
	_, mockOtpService, _, _, _, mockLogger, service := setupAuthTest(t)

	tests := []struct {
		name    string
		email   string
		userKey string
		mock    func()
		want    models.Otp
		wantErr bool
	}{
		{
			name:    "Success",
			email:   "test@example.com",
			userKey: "testUser",
			mock: func() {
				mockOtpService.EXPECT().
					RequestOTP("test@example.com", "testUser").
					Return(models.Otp{Email: "test@example.com"}, nil)
			},
			want:    models.Otp{Email: "test@example.com"},
			wantErr: false,
		},
		{
			name:    "OTP Service Error",
			email:   "test1@example.com",
			userKey: "testUser",
			mock: func() {
				mockOtpService.EXPECT().
					RequestOTP("test1@example.com", "testUser").
					Return(models.Otp{}, errs.InternalServerError(errors.New("error")))

				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

			},
			want:    models.Otp{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.RequestAuth(tt.email, tt.userKey)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVerifyAuth(t *testing.T) {
	mockAuthRepo, mockOtpService, mockAnalyticsService, _, mockJwt, _, service := setupAuthTest(t)

	tests := []struct {
		name      string
		email     string
		code      string
		requestID string
		mock      func()
		want      string
		wantErr   bool
	}{
		{
			name:      "Success - New User",
			email:     "new@example.com",
			code:      "123456",
			requestID: "req123",
			mock: func() {
				// Mock OTP verification
				mockOtpService.EXPECT().
					VerifyOTP("new@example.com", "123456", "req123").
					Return(models.Otp{Email: "new@example.com", Name: "New User"}, nil)

				// Mock user lookup by email - not found case
				mockAuthRepo.EXPECT().
					FindByEmail("new@example.com").
					Return(nil, gorm.ErrRecordNotFound)

				// Mock user lookup by handle - not found case
				mockAuthRepo.EXPECT().
					FindByHandle(mock.AnythingOfType("string")).
					Return(nil, gorm.ErrRecordNotFound)

				// Mock user creation
				mockAuthRepo.EXPECT().
					Save(mock.AnythingOfType("*models.User")).
					Return(nil)

				// Mock analytics
				mockAnalytics := &models.PlatformAnalytics{}
				mockAnalyticsService.EXPECT().
					GetCurrentData().
					Return(mockAnalytics)

				// Mock JWT generation
				mockJwt.EXPECT().
					GenerateToken(mock.AnythingOfType("uint"), "new@example.com", mock.AnythingOfType("string")).
					Return("token", nil)
			},
			want:    "token",
			wantErr: false,
		},
		// ...you can add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.VerifyAuth(tt.email, tt.code, tt.requestID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateUser(t *testing.T) {
	mockAuthRepo, _, mockAnalyticsService, _, _, _, service := setupAuthTest(t)

	tests := []struct {
		name    string
		user    models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Success - User Already Exists",
			user: models.User{Email: "test@example.com", Handle: "test"},
			mock: func() {
				mockAuthRepo.EXPECT().
					FindByHandle("test").
					Return(nil, errors.New("not found"))

				mockAuthRepo.EXPECT().
					Save(mock.AnythingOfType("*models.User")).
					Return(nil)

				mockAnalyticsService.EXPECT().
					GetCurrentData().
					Return(&models.PlatformAnalytics{})
			},
			wantErr: false,
		},
		{
			name: "Success new User",
			user: models.User{Email: "test@example.com", Handle: "test"},
			mock: func() {
				mockAuthRepo.EXPECT().
					FindByHandle("test").
					Return(&models.User{}, gorm.ErrRecordNotFound)
				mockAnalyticsService.EXPECT().
					GetCurrentData().
					Return(&models.PlatformAnalytics{})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetUserByHandle(t *testing.T) {
	mockAuthRepo, _, _, _, _, mockLogger, service := setupAuthTest(t)

	tests := []struct {
		name    string
		handle  string
		mock    func()
		want    models.User
		wantErr bool
	}{
		{
			name:   "Success",
			handle: "test",
			mock: func() {
				mockAuthRepo.EXPECT().
					FindByHandle("test").
					Return(&models.User{Handle: "test"}, nil)
			},
			want:    models.User{Handle: "test"},
			wantErr: false,
		},
		{
			name:   "User Not Found",
			handle: "nonexistent",
			mock: func() {
				mockAuthRepo.EXPECT().
					FindByHandle("nonexistent").
					Return(nil, errors.New("not found"))
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return()

			},
			want:    models.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetUserByHandle(tt.handle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
