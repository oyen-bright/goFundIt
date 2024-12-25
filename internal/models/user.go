package models

import (
	"time"

	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/utils"
	"gorm.io/gorm"
)

//TODO implement for mobile
// UserFCMToken{
//         Token:      req.Token,
//         UserHandle: userHandle,
//         Platform:   req.Platform,
//         CreatedAt:  time.Now(),
//     }

type User struct {
	ID            uint          `gorm:"primaryKey" json:"-"`
	Email         string        `gorm:"uniqueIndex;not null" encrypt:"true" binding:"required,email,lowercase"  validate:"email"`
	Handle        string        `gorm:"uniqueIndex;not null"`
	Name          *string       `gorm:"" encrypt:"true"`
	Verified      bool          `gorm:"not null" json:"-"`
	FCMToken      *string       `gorm:"" json:"-"`
	Contributions []Contributor `gorm:"foreignKey:Email;references:Email" json:"-"`
	CreatedAt     time.Time     `gorm:"not null" json:"-"`
	UpdatedAt     time.Time     `json:"-"`
}

func NewUser(name, email string, verified bool) *User {
	return &User{
		Name:     &name,
		Email:    email,
		Handle:   getHandle(email),
		Verified: verified,
	}
}

func (u *User) CanContributeToACampaign() bool {
	return len(u.Contributions) == 0

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

func (u *User) UpdateFCMToken(token string) {
	u.FCMToken = &token
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
