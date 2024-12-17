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

func (s *authService) RequestAuth(email, name string) (models.Otp, error) {
	otp, err := s.otpService.RequestOTP(email, name)
	if err != nil {
		return models.Otp{}, errs.InternalServerError(err).Log(s.logger)
	}
	return otp, nil
}

func (s *authService) VerifyAuth(email, code, requestID string) (string, error) {
	otp, err := s.otpService.VerifyOTP(email, code, requestID)
	if err != nil {
		return "", errs.BadRequest("Invalid OTP", err).Log(s.logger)
	}

	// Create User
	newUser := models.NewUser(otp.Name, otp.Email, true)
	if err := s.CreateUser(*newUser); err != nil {
		return "", err
	}

	// Generate Authentication token
	token, err := s.GenerateToken(*newUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

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

func (s *authService) FindNonExistingUsers(users []models.User) ([]models.User, error) {
	nonExistingUsers, err := s.authRepo.FindNonExistingUsers(users)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(s.logger)
	}
	return nonExistingUsers, nil
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

func (s *authService) GenerateToken(u models.User) (string, error) {
	token, err := s.jwt.GenerateToken(u.ID, u.Email, u.Handle)
	if err != nil {
		return "", errs.InternalServerError(err).Log(s.logger)
	}

	return token, nil
}
