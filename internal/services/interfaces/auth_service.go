package interfaces

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
)

type AuthService interface {
	RequestAuth(string, string) (models.Otp, error)
	VerifyAuth(string, string, string) (string, error)
	CreateUser(models.User) error
	CreateUsers(users []models.User) ([]models.User, error)
	FindExistingAndNonExistingUsers(email []string) (existing []models.User, nonExisting []string, err error)
	UpdateUser(models.User) error
	GetUserByHandle(handle string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	// DeleteUser(handle string) error
	GenerateToken(models.User) (string, error)
	GetAllUser() ([]models.User, error)
	SaveFCMToken(handle string, token string) error
	RemoveFCMToken(handle string, token string) error
	GetUsersByCreatedDateRange(from, to time.Time) ([]models.User, error)
}
