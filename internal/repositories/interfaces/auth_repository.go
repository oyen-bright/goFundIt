package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type AuthRepository interface {
	Save(user *models.User) error
	CreateMultiple(users []models.User) ([]models.User, error)
	Delete(user *models.User) error

	FindByHandle(handle string) (*models.User, error)
	FindNonExistingUsers(users []models.User) ([]models.User, error)
}
