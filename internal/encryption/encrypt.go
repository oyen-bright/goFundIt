package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(keys []string, data Data) (string, error) {
	//Encrypt data with the last key in the slice
	key := keys[len(keys)-1]
	encryptionKey, err := data.generateEncryptionKey(key)
	if err != nil {
		return "", err
	}
	encryptedData, err := encryptData(data.Data, encryptionKey)
	return encryptedData, err
}

func encryptData(data string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil

}
