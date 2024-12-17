package models

import (
	"time"

	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	ID            uint          `gorm:"primaryKey"`
	Email         string        `gorm:"uniqueIndex;not null" encrypt:"true" binding:"required,email,lowercase"  validate:"email"`
	Handle        string        `gorm:"uniqueIndex;not null"`
	Name          string        `gorm:"not null" encrypt:"true"`
	Verified      bool          `gorm:"not null"`
	Contributions []Contributor `gorm:"foreignKey:UserEmail;references:Email"`
	CreatedAt     time.Time     `gorm:"not null"`
	UpdatedAt     time.Time
}

func NewUser(name, email string, verified bool) *User {
	return &User{
		Name:     name,
		Email:    email,
		Handle:   getHandle(email),
		Verified: verified,
	}
}
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	//TODO:failing because data Already encrypted
	// validate := validator.New()
	// if err = validate.Struct(u); err != nil {
	// 	return err
	// }
	return nil
}

func getHandle(email string) string {
	return utils.GenerateRandomString(string(email[:2]), 0)
}
func (u *User) IsVerified() bool {
	return u.Verified
}

func (u *User) Verify() {
	u.Verified = true
}

func (u *User) Encrypt(e encryption.Encryptor) error {
	//TODO:disabled for now for faster dev
	return nil
	// var err error

	// encrypted, err := e.EncryptStruct(u, u.Email)
	// if err != nil {
	// 	return err
	// }

	// if user, ok := encrypted.(*User); ok {
	// 	*u = *user
	// }
	// return err

}

func (u *User) Decrypt(e encryption.Encryptor, key string) error {
	//TODO:disabled for now for faster dev
	return nil
	// var err error
	// encrypted, err := e.DecryptStruct(u, key)
	// if err != nil {
	// 	return err
	// }

	// if user, ok := encrypted.(*User); ok {
	// 	*u = *user
	// }
	// return err
}
