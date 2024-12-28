package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type AuthService interface {
	RequestAuth(string, string) (models.Otp, error)
	VerifyAuth(string, string, string) (string, error)

	CreateUser(models.User) error
	CreateUsers(users []models.User) ([]models.User, error)
	UpdateUser(models.User) error

	GetUserByHandle(handle string) (models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	FindExistingAndNonExistingUsers(email []string) (existing []models.User, nonExisting []string, err error)
	FindUserByEmail(email string) (*models.User, error)

	GenerateToken(models.User) (string, error)
	SaveFCMToken(handle string, token string) error
	RemoveFCMToken(handle string, token string) error

	GetAllUser() ([]models.User, error)
}
