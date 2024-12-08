package otp

import (
	"time"

	"math/rand"
	"strings"

	"github.com/oyen-bright/goFundIt/pkg/encryption"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyz"
const ExpireTime = 5 // minutes
const OTP_LENGTH = 6

type Otp struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null" encrypt:"true" binding:"required"`
	Code      string    `gorm:"not null"`
	Name      string    `encrypt:"true"`
	RequestId string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}

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
	encrypted, err := e.EncryptStruct(o, o.Email)
	if err != nil {
		return err
	}
	if otp, ok := encrypted.(*Otp); ok {
		*o = *otp
	}
	return err

}

func (o *Otp) Decrypt(e encryption.Encryptor, key string) error {
	encrypted, err := e.DecryptStruct(o, key)
	if err != nil {
		return err
	}
	if otp, ok := encrypted.(*Otp); ok {
		*o = *otp
	}
	return err
}

func New(email string) *Otp {

	code := GenerateOTP(0)
	expireAt := GenerateExpireDate()
	requestId := GenerateRequestId()
	return &Otp{
		Code:      code,
		Email:     strings.ToLower(email),
		CreatedAt: time.Now(),
		ExpiresAt: expireAt,
		RequestId: requestId,
	}
}

// GenerateOTP generates a random OTP
// if length is not provided, the default length is 6
func GenerateOTP(length int) string {

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

// generateExpireDate generates the time the OTP will expire
// the default expire time is 5 minutes
func GenerateExpireDate() time.Time {
	return time.Now().Add(time.Minute * ExpireTime)
}

func GenerateRequestId() string {

	numbers := make([]byte, 5)
	for i := range numbers {
		numbers[i] = '0' + byte(rand.Intn(10))
	}
	return string(numbers)
}
