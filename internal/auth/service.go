package auth

import (
	"github.com/oyen-bright/goFundIt/internal/utils/jwt"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
)

type AuthService interface {
	CreateUser(User) error
	UpdateUser(User) error
	GetUserByHandle(string, string) (User, error)
	DeleteUser(string) error
	GenerateToken(User) (string, error)
}

type authService struct {
	authRepo  AuthRepository
	encryptor encryptor.Encryptor
	jwt       jwt.Jwt
}

func Service(authRepo AuthRepository, encryptor encryptor.Encryptor, jwt jwt.Jwt) AuthService {
	return &authService{
		authRepo:  authRepo,
		encryptor: encryptor,
		jwt:       jwt,
	}
}

func (s *authService) CreateUser(u User) error {

	err := u.Encrypt(s.encryptor)

	if err != nil {
		return err
	}

	return s.authRepo.Save(&u)

}

func (s *authService) UpdateUser(u User) error {

	err := u.Encrypt(s.encryptor)

	if err != nil {
		return err
	}

	return s.authRepo.Save(&u)

}

func (s *authService) GetUserByHandle(handle, key string) (User, error) {

	var user *User

	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		return User{}, err
	}

	if err = user.Decrypt(s.encryptor, key); err != nil {
		return User{}, err
	}

	return *user, nil

}

func (s *authService) DeleteUser(handle string) error {
	user, err := s.authRepo.FindByHandle(handle)
	if err != nil {
		return err
	}
	return s.authRepo.Delete(user)

}

func (s *authService) GenerateToken(u User) (string, error) {

	return s.jwt.GenerateToken(u.ID, u.Email, u.Handle)
}
