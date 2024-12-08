package auth

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/utils"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null" encrypt:"true" binding:"required"`
	Handle    string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null" encrypt:"true"`
	Verified  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
}

func (u *User) IsVerified() bool {
	return u.Verified
}

func (u *User) Verify() {
	u.Verified = true
}

func (u *User) Encrypt(e encryption.Encryptor) error {
	var err error
	encrypted, err := e.EncryptStruct(u, u.Email)
	if err != nil {
		return err
	}
	if user, ok := encrypted.(*User); ok {
		*u = *user
	}
	return err
}

func (u *User) Decrypt(e encryption.Encryptor, key string) error {
	var err error
	encrypted, err := e.DecryptStruct(u, key)
	if err != nil {
		return err
	}
	if user, ok := encrypted.(*User); ok {
		*u = *user
	}
	return err
}

func New(name, email string, verified bool) *User {
	return &User{
		Name:     name,
		Email:    email,
		Handle:   getHandle(email),
		Verified: verified,
	}
}

func getHandle(email string) string {
	return utils.GenerateRandomString(string(email[:2]), 0)
}
