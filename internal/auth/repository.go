package auth

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository interface {
	save(user *User) error
	createMultiple(users []User) ([]User, error)
	delete(user *User) error

	FindByHandle(handle string) (*User, error)
	FindNonExistingUsers(users []User) ([]User, error)
}

type authRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// ----------------------------------------------------------------------
func (r *authRepository) save(user *User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "verified"}),
	}).Create(user).Error
}

func (r *authRepository) createMultiple(users []User) ([]User, error) {
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
func (r *authRepository) delete(user *User) error {
	return r.db.Create(user).Error
}

// ---------------------------------------------------------------------

func (r *authRepository) FindByHandle(handle string) (*User, error) {
	var user User
	err := r.db.Preload("Contributions").Where("handle = ?", handle).First(&user).Error

	log.Println(user)
	log.Println("xxxx")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindNonExistingUsers(users []User) ([]User, error) {
	var nonExistingUsers []User
	for _, user := range users {
		var existingUser User
		err := r.db.Where("email = ?", user.Email).First(&existingUser).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			nonExistingUsers = append(nonExistingUsers, user)
		} else if err != nil {
			return nil, err
		}
	}
	return nonExistingUsers, nil
}
