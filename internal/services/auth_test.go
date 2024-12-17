// internal/services/auth_test.go

package services_test

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/mocks"
	"github.com/oyen-bright/goFundIt/internal/services"
	serviceMocks "github.com/oyen-bright/goFundIt/internal/services/mocks"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type serviceDependencies struct {
	authRepo   *mocks.AuthRepository
	otpService *serviceMocks.OTPService
	encryptor  encryption.Encryptor
	jwt        jwt.Jwt
	logger     logger.Logger
}

func setupTest(t *testing.T) serviceDependencies {
	return serviceDependencies{
		authRepo:   mocks.NewAuthRepository(t),
		otpService: serviceMocks.NewOTPService(t),
		encryptor:  *encryption.New([]string{"test-secret"}),
		jwt:        jwt.New("test-secret"),
		logger:     logger.New(),
	}
}

func TestAuthService_RequestAuth(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name     string
		email    string
		userName string
		mockFn   func()
		wantErr  bool
	}{
		{
			name:     "successful request",
			email:    "test@example.com",
			userName: "Test User",
			mockFn: func() {
				deps.otpService.On("RequestOTP", "test@example.com", "Test User").
					Return(models.Otp{Email: "test@example.com"}, nil)
			},
			wantErr: false,
		},
		{
			name:     "otp service error",
			email:    "error@example.com",
			userName: "Error User",
			mockFn: func() {
				deps.otpService.On("RequestOTP", "error@example.com", "Error User").
					Return(models.Otp{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			otp, err := authService.RequestAuth(tt.email, tt.userName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.email, otp.Email)
			deps.otpService.AssertExpectations(t)
		})
	}
}

func TestAuthService_VerifyAuth(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name      string
		email     string
		code      string
		requestID string
		mockFn    func()
		wantErr   bool
	}{
		{
			name:      "successful verification",
			email:     "test@example.com",
			code:      "123456",
			requestID: "req-id",
			mockFn: func() {
				deps.otpService.On("VerifyOTP", "test@example.com", "123456", "req-id").
					Return(models.Otp{Email: "test@example.com", Name: "Test User"}, nil)
				deps.authRepo.On("FindByHandle", mock.AnythingOfType("string")).
					Return(nil, gorm.ErrRecordNotFound)
				deps.authRepo.On("Save", mock.AnythingOfType("*models.User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "invalid otp",
			email:     "test@example.com",
			code:      "invalid",
			requestID: "req-id",
			mockFn: func() {
				deps.otpService.On("VerifyOTP", "test@example.com", "invalid", "req-id").
					Return(models.Otp{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			token, err := authService.VerifyAuth(tt.email, tt.code, tt.requestID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			deps.otpService.AssertExpectations(t)
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_CreateUser(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		user    models.User
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful create",
			user: models.User{Email: "new@example.com", Name: "New User"},
			mockFn: func() {
				deps.authRepo.On("FindByHandle", mock.AnythingOfType("string")).
					Return(nil, gorm.ErrRecordNotFound)
				deps.authRepo.On("Save", mock.AnythingOfType("*models.User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			user: models.User{Email: "existing@example.com", Name: "Existing User"},
			mockFn: func() {
				deps.authRepo.On("FindByHandle", mock.AnythingOfType("string")).
					Return(&models.User{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := authService.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GetUserByHandle(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		handle  string
		mockFn  func()
		want    models.User
		wantErr bool
	}{
		{
			name:   "user found",
			handle: "test-handle",
			mockFn: func() {
				deps.authRepo.On("FindByHandle", "test-handle").
					Return(&models.User{
						Handle: "test-handle",
						Email:  "test@example.com",
						Name:   "Test User",
					}, nil)
			},
			want: models.User{
				Handle: "test-handle",
				Email:  "test@example.com",
				Name:   "Test User",
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			handle: "non-existent",
			mockFn: func() {
				deps.authRepo.On("FindByHandle", "non-existent").
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			user, err := authService.GetUserByHandle(tt.handle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, user)
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_DeleteUser(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		handle  string
		mockFn  func()
		wantErr bool
	}{
		{
			name:   "successful delete",
			handle: "test-handle",
			mockFn: func() {
				deps.authRepo.On("FindByHandle", "test-handle").
					Return(&models.User{Handle: "test-handle"}, nil)
				deps.authRepo.On("Delete", mock.AnythingOfType("*models.User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			handle: "non-existent",
			mockFn: func() {
				deps.authRepo.On("FindByHandle", "non-existent").
					Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			err := authService.DeleteUser(tt.handle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_CreateUsers(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		users   []models.User
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful batch create",
			users: []models.User{
				{Email: "user1@example.com", Name: "User 1"},
				{Email: "user2@example.com", Name: "User 2"},
			},
			mockFn: func() {
				deps.authRepo.On("CreateMultiple", mock.AnythingOfType("[]models.User")).
					Return([]models.User{
						{Email: "user1@example.com", Name: "User 1"},
						{Email: "user2@example.com", Name: "User 2"},
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "creation error",
			users: []models.User{
				{Email: "error@example.com", Name: "Error User"},
			},
			mockFn: func() {
				deps.authRepo.On("CreateMultiple", mock.AnythingOfType("[]models.User")).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			users, err := authService.CreateUsers(tt.users)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, users, len(tt.users))
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_FindNonExistingUsers(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		users   []models.User
		mockFn  func()
		want    []models.User
		wantErr bool
	}{
		{
			name: "find non-existing users",
			users: []models.User{
				{Email: "new1@example.com", Name: "New User 1"},
				{Email: "new2@example.com", Name: "New User 2"},
			},
			mockFn: func() {
				deps.authRepo.On("FindNonExistingUsers", mock.AnythingOfType("[]models.User")).
					Return([]models.User{
						{Email: "new1@example.com", Name: "New User 1"},
					}, nil)
			},
			want: []models.User{
				{Email: "new1@example.com", Name: "New User 1"},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			users: []models.User{
				{Email: "error@example.com", Name: "Error User"},
			},
			mockFn: func() {
				deps.authRepo.On("FindNonExistingUsers", mock.AnythingOfType("[]models.User")).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			users, err := authService.FindNonExistingUsers(tt.users)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, users)
			deps.authRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	deps := setupTest(t)
	authService := services.NewAuthService(deps.authRepo, deps.otpService, deps.encryptor, deps.jwt, deps.logger)

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "successful token generation",
			user: models.User{
				ID:     1,
				Email:  "test@example.com",
				Handle: "test-handle",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.GenerateToken(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		})
	}
}
