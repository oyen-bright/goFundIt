package encryption

import (
	"errors"
	"reflect"
)

type Encryptor struct {
	Keys []string
}

func New(key []string) *Encryptor {
	return &Encryptor{Keys: key}
}

func (e *Encryptor) Encrypt(data Data) (string, error) {
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

func (e *Encryptor) EncryptStruct(data interface{}, key string) (interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Check if the field has the "encrypt" tag set to "true"
		if tag, ok := fieldType.Tag.Lookup("encrypt"); ok && tag == "true" {

			//prepare the data to be encrypted
			modelData := Data{
				Key:  key,
				Data: field.Interface().(string),
			}

			// Encrypt the field value
			encryptedData, err := e.Encrypt(modelData)
			if err != nil {
				return data, err
			}

			// Ensure the field is settable
			if field.CanSet() && field.Kind() == reflect.String {
				field.SetString(encryptedData)
			} else {
				return data, errors.New("field is not settable or not a string")
			}
		}
	}

	return data, nil
}

func (e *Encryptor) Decrypt(data Data) (string, error) {

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

func (e *Encryptor) DecryptStruct(data interface{}, key string) (interface{}, error) {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return data, errors.New("input data is not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Check if the field has the "encrypt" tag set to "true"
		if tag, ok := fieldType.Tag.Lookup("encrypt"); ok && tag == "true" {
			modelData := Data{
				Key:  key,
				Data: field.Interface().(string),
			}

			// Try to decrypt using available keys
			var decryptErr error
			for _, sKey := range e.Keys {
				encryptionKey, err := modelData.GenerateEncryptionKey(sKey)
				if err != nil {
					decryptErr = err
					continue
				}

				decryptedData, err := decryptData(modelData.Data, encryptionKey)
				if err == nil {
					if field.CanSet() && field.Kind() == reflect.String {
						field.SetString(decryptedData)
						decryptErr = nil
						break
					} else {
						return data, errors.New("field is not settable or not a string")
					}
				} else {
					decryptErr = err
				}
			}

			if decryptErr != nil {
				return data, decryptErr
			}
		}
	}

	return data, nil
}
