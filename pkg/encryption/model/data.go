package model

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/oyen-bright/goFundIt/internal/utils/validator"
)

type Data struct {
	Data  string
	Email string
	Id    string
}

func (e Data) GenerateEncryptionKey(secretKey string) (string, error) {
	// Validate the secret key

	if secretKey == "" {

		return "", errors.New("invalid secret key")
	}

	// Validate the email
	isValidEmail := validator.ValidateEmail(e.Email)
	if !isValidEmail {
		return "", errors.New("invalid email")
	}

	// Generate a SHA-256 hash of the email and secret key
	hash := sha256.Sum256([]byte(e.Email + secretKey))

	// Encode the first 24 bytes of the hash to a base64 string
	return base64.StdEncoding.EncodeToString(hash[:24]), nil
}
