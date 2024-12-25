package services

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

// TODO(oyen-bright): Split this service into separate AuthService and UserService
// AuthService should handle:
// - Authentication (RequestAuth, VerifyAuth, GenerateToken)
// - Token management
// UserService should handle:
// - User CRUD operations
// - User queries
// Created: 2024-12-18
// Priority: Medium

// authService handles both authentication and user management
// This service currently has two responsibilities and should be split
// into separate services for better separation of concerns
type authService struct {
	authRepo   repositories.AuthRepository
	otpService services.OTPService
	encryptor  encryptor.Encryptor
	jwt        jwt.Jwt
	logger     logger.Logger
}

func NewAuthService(
	authRepo repositories.AuthRepository,
	otpService services.OTPService,
	encryptor encryptor.Encryptor,
	jwt jwt.Jwt,
	logger logger.Logger,
) services.AuthService {
	return &authService{
		authRepo:   authRepo,
		otpService: otpService,
		encryptor:  encryptor,
		jwt:        jwt,
		logger:     logger,
	}
}

// Authentication Methods
func (s *authService) RequestAuth(email, name string) (models.Otp, error) {
	otp, err := s.otpService.RequestOTP(email, name)
	if err != nil {
		// return models.Otp{}, errs.InternalServerError(err).Log(s.logger)
	}
	return otp, nil
}

func (s *authService) VerifyAuth(email, code, requestID string) (string, error) {
	// Verify OTP
	otp, err := s.otpService.VerifyOTP(email, code, requestID)
	if err != nil {
		return "", errs.BadRequest("Invalid OTP", err).Log(s.logger)
	}

	// Get or create user
	user, err := s.getOrCreateUser(otp)
	if err != nil {
		return "", err
	}

	// Generate Authentication token
	return s.GenerateToken(*user)
}

func (s *authService) GenerateToken(u models.User) (string, error) {
	token, err := s.jwt.GenerateToken(u.ID, u.Email, u.Handle)
	if err != nil {
		return "", errs.InternalServerError(err).Log(s.logger)
	}
	return token, nil
}

// User Management Methods
func (s *authService) CreateUser(u models.User) error {
	// Check if user already exists
	_, err := s.authRepo.FindByHandle(u.Handle)
	if err == nil {
		return errs.BadRequest("User already exists", nil)
	}

	if err := u.Encrypt(s.encryptor); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	if err := s.authRepo.Save(&u); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil
}

func (s *authService) CreateUsers(users []models.User) ([]models.User, error) {
	for i := range users {
		if err := users[i].Encrypt(s.encryptor); err != nil {
			return nil, errs.InternalServerError(err).Log(s.logger)
		}
	}

	createdUsers, err := s.authRepo.CreateMultiple(users)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}

	return createdUsers, nil
}

func (s *authService) UpdateUser(u models.User) error {
	if err := u.Encrypt(s.encryptor); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	if err := s.authRepo.Save(&u); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil
}

func (s *authService) DeleteUser(handle string) error {
	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.BadRequest("User not found", nil)
		}
		return errs.InternalServerError(err).Log(s.logger)
	}

	if err := s.authRepo.Delete(user); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil
}

// User Query Methods
func (s *authService) GetUserByHandle(handle string) (models.User, error) {
	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return models.User{}, errs.BadRequest("User not found", nil)
		}
		return models.User{}, errs.InternalServerError(err).Log(s.logger)
	}

	return *user, nil
}

func (s *authService) GetUserByEmail(email string) (models.User, error) {
	user, err := s.authRepo.FindByEmail(email)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return models.User{}, errs.BadRequest("User not found", nil)
		}
		return models.User{}, errs.InternalServerError(err).Log(s.logger)
	}

	return *user, nil
}

func (s *authService) FindExistingAndNonExistingUsers(emails []string) (existing []models.User, nonExisting []string, err error) {
	existing, nonExisting, err = s.authRepo.FindExistingAndNonExistingUsers(emails)
	if err != nil {
		return nil, nil, errs.InternalServerError(err).Log(s.logger)
	}
	return existing, nonExisting, nil
}

// GetAllUser retrieves all users from the repository
func (s *authService) GetAllUser() ([]models.User, error) {
	users, err := s.authRepo.GetAll()
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return users, nil
}

// Helper methods
func (s *authService) getOrCreateUser(otp models.Otp) (*models.User, error) {
	user, err := s.authRepo.FindByEmail(otp.Email)
	if err == nil {
		return user, nil
	}

	if !database.Error(err).IsNotfound() {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}

	// Create new user
	newUser := models.NewUser(otp.Name, otp.Email, true)
	if err := s.CreateUser(*newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// UpdateFCMToken updates the FCM token for a user
func (s *authService) SaveFCMToken(handle string, token string) error {
	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.BadRequest("User not found", nil)
		}
		return errs.InternalServerError(err).Log(s.logger)
	}

	user.FCMToken = &token
	if err := s.authRepo.Save(user); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil
}

// RemoveFCMToken removes the FCM token for a user if it matches the provided token
func (s *authService) RemoveFCMToken(handle string, token string) error {
	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return errs.BadRequest("User not found", nil)
		}
		return errs.InternalServerError(err).Log(s.logger)
	}

	if user.FCMToken == nil || *user.FCMToken != token {
		return nil
	}

	user.FCMToken = nil
	if err := s.authRepo.Save(user); err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil
}
