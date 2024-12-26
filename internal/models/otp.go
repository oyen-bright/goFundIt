package models

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oyen-bright/goFundIt/pkg/encryption"
)

const (
	CHARSET    = "abcdefghijklmnopqrstuvwxyz"
	ExpireTime = 5 // minutes
	OTP_LENGTH = 6
)

type Otp struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null" encrypt:"true" binding:"required,email"`
	Code      string    `gorm:"not null"`
	Name      string    `encrypt:"true"`
	RequestId string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

// Constructor
func NewOTP(email string) *Otp {
	return &Otp{
		Email:     strings.ToLower(email),
		Code:      generateOTP(0),
		RequestId: generateRequestId(),
		CreatedAt: time.Now(),
		ExpiresAt: generateExpireDate(),
	}
}

// Methods
func (o *Otp) IsExpired() bool {
	return o.ExpiresAt.Before(time.Now())
}

func (o *Otp) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"request_id": o.RequestId,
		//TODO:remove
		"code": o.Code,
	}
}

func (o *Otp) Encrypt(e encryption.Encryptor) error {
	//TODO:disabled for now for faster dev
	return nil
	// encrypted, err := e.EncryptStruct(o, o.Email)
	// if err != nil {
	// 	return err
	// }
	// if otp, ok := encrypted.(*Otp); ok {
	// 	*o = *otp
	// }
	// return err
}

func (o *Otp) Decrypt(e encryption.Encryptor, key string) error {
	//TODO:disabled for now for faster dev
	return nil
	// encrypted, err := e.DecryptStruct(o, key)
	// if err != nil {
	// 	return err
	// }
	// if otp, ok := encrypted.(*Otp); ok {
	// 	*o = *otp
	// }
	// return err
}

// Helper functions ------------------------------------

func generateOTP(length int) string {
	otpLength := length
	if length == 0 {
		otpLength = OTP_LENGTH
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))

	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = CHARSET[rand.Intn(len(CHARSET))]
	}
	return strings.ToUpper(string(otp))
}

func generateExpireDate() time.Time {
	return time.Now().Add(time.Minute * ExpireTime)
}

func generateRequestId() string {
	numbers := make([]byte, 5)
	for i := range numbers {
		numbers[i] = '0' + byte(rand.Intn(10))
	}
	return string(numbers)
}
