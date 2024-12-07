package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
)

// TODO: reduce the size of the encrypted data
// TODO: consider using a different encryption algorithm or different nonce for each data
func encryptData(data string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceHash := sha256.Sum256([]byte(data))
	nonce := nonceHash[:aesGCM.NonceSize()]
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil

}
