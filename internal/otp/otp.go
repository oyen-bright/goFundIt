package otp

import (
	"math/rand"
	"strings"
	"time"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyz"
const ExpireTime = 5 // minutes
const OTP_LENGTH = 6

func New(email string) Otp {

	code := GenerateOTP(0)
	expireAt := GenerateExpireDate()
	return Otp{
		Code:      code,
		Email:     strings.ToLower(email),
		CreatedAt: time.Now(),
		ExpiresAt: expireAt,
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
