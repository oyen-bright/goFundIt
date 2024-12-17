package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &authRepository{db: db}
}

// ----------------------------------------------------------------------
func (r *authRepository) Save(user *models.User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "verified"}),
	}).Create(user).Error
}

func (r *authRepository) CreateMultiple(users []models.User) ([]models.User, error) {
	batchSize := 100
	if len(users) < batchSize {
		batchSize = len(users)
	}

	result := r.db.CreateInBatches(&users, batchSize)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
func (r *authRepository) Delete(user *models.User) error {
	return r.db.Create(user).Error
}

// ---------------------------------------------------------------------

func (r *authRepository) FindByHandle(handle string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Contributions").Where("handle = ?", handle).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindNonExistingUsers(users []models.User) ([]models.User, error) {
	var nonExistingUsers []models.User
	for _, user := range users {
		var existingUser models.User
		err := r.db.Where("email = ?", user.Email).First(&existingUser).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			nonExistingUsers = append(nonExistingUsers, user)
		} else if err != nil {
			return nil, err
		}
	}
	return nonExistingUsers, nil
}
