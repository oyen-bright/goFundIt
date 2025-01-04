package utils

import (
	"math/rand"
	"strings"
	"time"
)

// GenerateRandomString generates a random alphanumeric string of the specified length.
//
//   - If the provided length is 0, it defaults to a predefined OTP_LENGTH.
//   - if prefix is not empty its is added in front ot the generated string
func GenerateRandomString(prefix string, length int) string {
	const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const ID_LENGTH = 6
	IDLength := length
	if length == 0 {
		IDLength = ID_LENGTH
	}
	IDLength -= (len(prefix))
	rand.New(rand.NewSource(time.Now().UnixNano()))

	otp := make([]byte, IDLength)
	for i := range otp {
		otp[i] = CHARSET[rand.Intn(len(CHARSET))]
	}

	return strings.ToUpper(prefix + string(otp))
}

// GenerateRandomNumber generates a random number of the specified length.

func GenerateRandomAlphaNumeric(prefix string, length int) string {
	const NUM_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))

	idLength := length - len(prefix)
	id := make([]byte, idLength)
	for i := range id {
		id[i] = NUM_CHARSET[rand.Intn(len(NUM_CHARSET))]
	}

	return strings.ToUpper(prefix + string(id))
}
