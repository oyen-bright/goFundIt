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

	nonceHash := sha256.Sum256([]byte(data))
	nonce := nonceHash[:aesGCM.NonceSize()]
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil

}

// func EncryptDataM(data, secretKey string, keys []string) (string, error) {
// 	if len(keys) == 0 {
// 		return "", errors.New("no key provided")
// 	}

// 	// Combine all ke + secret key into a single hash
// 	masterHash := deriveMasterHash(keys, secretKey)

// 	// XOR the data with the hash
// 	encrypted := xorData([]byte(data), []byte(masterHash))

// 	return hex.EncodeToString(encrypted), nil
// }

// func deriveMasterHash(keys []string, secretKey string) string {
// 	hash := sha256.New()
// 	for _, key := range keys {
// 		hash.Write([]byte(key))
// 	}
// 	hash.Write([]byte(secretKey))
// 	return hex.EncodeToString(hash.Sum(nil))
// }

// func xorData(data, key []byte) []byte {
// 	keyLen := len(key)
// 	result := make([]byte, len(data))
// 	for i := 0; i < len(data); i++ {
// 		result[i] = data[i] ^ key[i%keyLen]
// 	}
// 	return result
// }
