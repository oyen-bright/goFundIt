package auth

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
	FindByHandle(handle string) (*User, error)
	Save(user *User) error
	Delete(user *User) error
}

type authRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByHandle(handle string) (*User, error) {
	var user User
	err := r.db.Where("handle = ?", handle).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Save(user *User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},                       // Specify the conflict target
		DoUpdates: clause.AssignmentColumns([]string{"name", "verified"}), // Columns to update on conflict
	}).Create(user).Error
}

func (r *authRepository) Delete(user *User) error {
	return r.db.Create(user).Error
}
