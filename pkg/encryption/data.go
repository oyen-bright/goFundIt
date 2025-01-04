package encryption

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

type Data struct {
	Data string
	Key  string
	Id   string
}

func (e Data) GenerateEncryptionKey(secretKey string) (string, error) {
	// Validate the secret key

	if secretKey == "" {

		return "", errors.New("invalid secret key")
	}

	// Generate a SHA-256 hash of the email and secret key
	hash := sha256.Sum256([]byte(strings.ToLower(e.Key) + secretKey))

	// Encode the first 24 bytes of the hash to a base64 string
	return base64.StdEncoding.EncodeToString(hash[:24]), nil
}

func NewData(data, key, id string) *Data {
	return &Data{
		Data: data,
		Key:  key,
		Id:   id,
	}
}
