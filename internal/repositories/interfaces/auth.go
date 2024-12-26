package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type AuthRepository interface {
	Save(user *models.User) error
	CreateMultiple(users []models.User) ([]models.User, error)
	Delete(user *models.User) error

	GetAll() ([]models.User, error)
	FindByHandle(handle string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindExistingAndNonExistingUsers(emails []string) (existing []models.User, nonExisting []string, err error)
}
