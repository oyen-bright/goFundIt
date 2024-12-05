package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func Decrypt(keys []string, data Data) (string, error) {
	var plaintext *string

	// Reverse the slice
	// Loop through the keys to find a match to decrypt the data
	for index, key := range keys {
		encryptionKey, err := data.generateEncryptionKey(key)
		if err != nil {
			return "", err
		}

		data, err := decryptData(data.Data, encryptionKey)

		if err == nil {
			plaintext = &data
			break
		}

		if index < len(keys)-1 {
			continue
		}
		return "", err

	}
	return *plaintext, nil

}

func decryptData(data string, key string) (string, error) {

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	data = string(decodedData)

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil

}
