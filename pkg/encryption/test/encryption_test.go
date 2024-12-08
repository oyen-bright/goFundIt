package encryption_test

import (
	"os"
	"testing"

	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/environment"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/encryption/model"
)

var appConfig *config.AppConfig
var encryptor *encryption.Encryptor

func mockConfig() *config.AppConfig {
	return &config.AppConfig{
		Environment:   environment.Development,
		Port:          ":8080",
		EmailProvider: providers.EmailSMTP,
		EncryptionKey: []string{"test-key"},
	}
}

// TestMain initializes the app configuration and sets it up for tests
func TestMain(m *testing.M) {

	appConfig = mockConfig()
	encryptor = encryption.New(appConfig.EncryptionKey)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		email       string
		expectError bool
	}{
		{
			name:        "Valid encryption and decryption",
			data:        "Sensitive campaign data",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "Invalid email",
			data:        "Sensitive campaign data",
			email:       "invalid-email",
			expectError: true,
		},
		{
			name:        "Empty data",
			data:        "",
			email:       "test@example.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the data
			encryptedData, err := encryptor.Encrypt(model.Data{Data: tt.data, Key: tt.email})
			if (err != nil) != tt.expectError {
				t.Fatalf("Encrypt() error = %v, expectError = %v", err, tt.expectError)
			}

			// Skip decryption test if encryption fails
			if tt.expectError {
				return
			}

			// Decrypt the data
			decryptedData, err := encryptor.Decrypt(model.Data{Data: encryptedData, Key: tt.email})
			if err != nil {
				t.Fatalf("Failed to decrypt data: %v", err)
			}

			// Check if the decrypted data matches the original data
			if decryptedData != tt.data {
				t.Errorf("Decrypted data does not match original data. Got: %s, Want: %s", decryptedData, tt.data)
			}
		})
	}
}

func TestMissingSecretKey(t *testing.T) {
	encryptor.Keys = []string{}

	data := "Sensitive campaign data"
	email := "test@example.com"

	_, err := encryptor.Encrypt(model.Data{Data: data, Key: email})
	if err == nil {
		t.Fatalf("Expected error for missing secret key, but got none")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	originalKey := appConfig.EncryptionKey
	wrongKey := []string{"wrong-secure-secret-key"}
	data := "Sensitive campaign data"
	email := "test@example.com"

	encryptor.Keys = originalKey

	// Encrypt with the original key
	encryptedData, err := encryptor.Encrypt(model.Data{Data: data, Key: email})
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	encryptor.Keys = wrongKey

	// Try to decrypt with the wrong key
	_, err = encryptor.Decrypt(model.Data{Data: encryptedData, Key: email})
	if err == nil {
		t.Fatalf("Expected error for decrypting with the wrong key, but got none")
	}
}
