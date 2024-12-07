package encryption

import (
	"errors"

	"github.com/oyen-bright/goFundIt/pkg/encryption/model"
)

type Encryptor struct {
	Keys []string
}

func New(key []string) *Encryptor {
	return &Encryptor{Keys: key}
}

func (e *Encryptor) Encrypt(data model.Data) (string, error) {
	// Check if the secret key is missing
	if len(e.Keys) == 0 {
		return "", errors.New("missing secret key")
	}

	//Encrypt data with the last key in the slice
	key := e.Keys[len(e.Keys)-1]
	encryptionKey, err := data.GenerateEncryptionKey(key)
	if err != nil {
		return "", err
	}
	return encryptData(data.Data, encryptionKey)
}

func (e *Encryptor) EncryptT(email, data string) (string, error) {

	modelData := model.Data{
		Email: email,
		Data:  data,
	}
	// Check if the secret key is missing
	if len(e.Keys) == 0 {
		return "", errors.New("missing secret key")
	}

	//Encrypt data with the last key in the slice
	key := e.Keys[len(e.Keys)-1]
	encryptionKey, err := modelData.GenerateEncryptionKey(key)
	if err != nil {
		return "", err
	}
	return encryptData(modelData.Data, encryptionKey)
}

func (e *Encryptor) Decrypt(data model.Data) (string, error) {

	var plaintext *string

	// Reverse the slice
	// Loop through the keys to find a match to decrypt the data
	for index, key := range e.Keys {
		encryptionKey, err := data.GenerateEncryptionKey(key)
		if err != nil {
			return "", err
		}

		data, err := decryptData(data.Data, encryptionKey)

		if err == nil {
			plaintext = &data
			break
		}

		if index < len(e.Keys)-1 {
			continue
		}
		return "", err

	}
	return *plaintext, nil
}
