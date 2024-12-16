package auth

import (
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type AuthService interface {
	RequestAuth(string, string) (otp.Otp, error)
	VerifyAuth(string, string, string) (string, error)
	CreateUser(User) error
	CreateUsers(users []User) ([]User, error)
	FindNonExistingUsers(users []User) ([]User, error)
	UpdateUser(User) error
	GetUserByHandle(handle string) (User, error)
	DeleteUser(string) error
	GenerateToken(User) (string, error)
}

type authService struct {
	authRepo   AuthRepository
	otpService otp.OTPService
	encryptor  encryptor.Encryptor
	jwt        jwt.Jwt
	logger     logger.Logger
}

func Service(authRepo AuthRepository, otpService otp.OTPService, encryptor encryptor.Encryptor, jwt jwt.Jwt, logger logger.Logger) AuthService {
	return &authService{
		authRepo:   authRepo,
		otpService: otpService,
		encryptor:  encryptor,
		jwt:        jwt,
		logger:     logger,
	}
}

func (s *authService) RequestAuth(email, name string) (otp.Otp, error) {
	return s.otpService.RequestOTP(email, name)
}

func (s *authService) VerifyAuth(email, code, requestID string) (string, error) {
	otp, err := s.otpService.VerifyOTP(email, code, requestID)
	if err != nil {

		return "", err
	}
	//Create User
	newUser := New(otp.Name, otp.Email, true)
	err = s.CreateUser(*newUser)
	if err != nil {
		return "", err

	}

	//Generate Authentication token
	token, err := s.GenerateToken(*newUser)
	if err != nil {
		return "", err

	}

	return token, nil

}

//---------------------------------------------------------------------------

func (s *authService) CreateUser(u User) error {

	err := u.Encrypt(s.encryptor)
	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	err = s.authRepo.save(&u)
	if err != nil {
		return errs.InternalServerError(err).Log(s.logger)
	}

	return nil

}

func (s *authService) CreateUsers(users []User) ([]User, error) {

	for _, u := range users {

		err := u.Encrypt(s.encryptor)
		if err != nil {
			return users, errs.InternalServerError(err).Log(s.logger)
		}
	}

	users, err := s.authRepo.createMultiple(users)
	if err != nil {
		return users, errs.InternalServerError(err).Log(s.logger)
	}

	return users, nil

}

func (s *authService) FindNonExistingUsers(users []User) ([]User, error) {

	users, err := s.authRepo.FindNonExistingUsers(users)
	if err != nil {
		return users, errs.InternalServerError(err).Log(s.logger)
	}
	return users, nil
}

func (s *authService) UpdateUser(u User) error {

	err := u.Encrypt(s.encryptor)
	if err != nil {
		return err
	}
	return s.authRepo.save(&u)
}

func (s *authService) GetUserByHandle(handle string) (User, error) {

	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {

		if errs.NewDB(err).IsNotfound() {
			return User{}, errs.BadRequest("User not found", nil)
		}

		return User{}, errs.InternalServerError(err).Log(s.logger)
	}

	// if err = user.Decrypt(s.encryptor, key); err != nil {
	// 	return User{}, errs.InternalServerError(err).Log(s.logger)
	// }

	return *user, nil
}

func (s *authService) DeleteUser(handle string) error {

	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		return err
	}

	return s.authRepo.delete(user)
}

func (s *authService) GenerateToken(u User) (string, error) {

	token, err := s.jwt.GenerateToken(u.ID, u.Email, u.Handle)
	if err != nil {
		return "", errs.InternalServerError(err).Log(s.logger)
	}

	return token, nil
}
