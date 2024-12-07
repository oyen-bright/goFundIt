package otp

import (
	"time"

	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/encryption/model"
)

type Otp struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

func (o *Otp) IsExpired() bool {
	return o.ExpiresAt.Before(time.Now())
}

func (o *Otp) Encrypt(e encryption.Encryptor) error {
	var err error
	o.Email, err = e.Encrypt(model.Data{
		Email: o.Email,
		Data:  o.Email,
	})
	return err
}

func (o *Otp) Decrypt(e encryption.Encryptor) error {
	var err error
	o.Email, err = e.Decrypt(model.Data{
		Email: o.Email,
		Data:  o.Email,
	})
	return err
}
