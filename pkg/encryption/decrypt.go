package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func decryptData(data string, key string) (string, error) {

	// Don't encrypt if the data is empty
	if data == "" {
		return "", nil
	}

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

// func DecryptDataM(encryptedData, secretKey, email string) (string, error) {
// 	// Decode the hex-encoded encrypted data
// 	encryptedBytes, err := hex.DecodeString(encryptedData)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Derive the hash from the given email + secret key
// 	emailHash := deriveHash(email, secretKey)

// 	// XOR the encrypted data with the hash to decrypt
// 	decrypted := xorData(encryptedBytes, []byte(emailHash))

// 	return string(decrypted), nil
// }

// // deriveHash generates a hash for a single email + the secret key
// func deriveHash(email, secretKey string) string {
// 	hash := sha256.New()
// 	hash.Write([]byte(email))
// 	hash.Write([]byte(secretKey))
// 	return hex.EncodeToString(hash.Sum(nil))
// }
