package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type AuthService interface {
	RequestAuth(string, string) (models.Otp, error)
	VerifyAuth(string, string, string) (string, error)
	CreateUser(models.User) error
	CreateUsers(users []models.User) ([]models.User, error)
	FindNonExistingUsers(users []models.User) ([]models.User, error)
	UpdateUser(models.User) error
	GetUserByHandle(handle string) (models.User, error)
	DeleteUser(handle string) error
	GenerateToken(models.User) (string, error)
}
