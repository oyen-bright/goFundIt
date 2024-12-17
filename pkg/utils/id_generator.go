package utils

import (
	"math/rand"
	"strings"
	"time"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const OTP_LENGTH = 6

// GenerateRandomString generates a random alphanumeric string of the specified length.
//
//   - If the provided length is 0, it defaults to a predefined OTP_LENGTH.
//   - if prefix is not empty its is added in front ot the generated string
func GenerateRandomString(prefix string, length int) string {

	otpLength := length
	if length == 0 {
		otpLength = OTP_LENGTH
	}
	otpLength -= (len(prefix))
	rand.New(rand.NewSource(time.Now().UnixNano()))

	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = CHARSET[rand.Intn(len(CHARSET))]
	}

	return strings.ToUpper(prefix + string(otp))
}
