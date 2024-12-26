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

func (r *authRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil

}

func (r *authRepository) FindByHandle(handle string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Contributions").Where("handle = ?", handle).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Contributions").Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindExistingAndNonExistingUsers(emails []string) (existing []models.User, nonExisting []string, err error) {
	// Find all existing users with their contributions in a single query
	var existingUsers []models.User
	if err := r.db.Preload("Contributions").Where("email IN ?", emails).Find(&existingUsers).Error; err != nil {
		return nil, nil, err
	}

	// Create a map of found emails for faster lookup
	foundEmails := make(map[string]bool)
	for _, user := range existingUsers {
		foundEmails[user.Email] = true
	}

	// Determine non-existing emails
	for _, email := range emails {
		if !foundEmails[email] {
			nonExisting = append(nonExisting, email)
		}
	}

	return existingUsers, nonExisting, nil
}
