package encryption

import (
	"os"
	"testing"
)

var (
	testEncryptor Encryptor
)

// TestMain initializes the app configuration and sets it up for tests
func TestMain(m *testing.M) {
	testEncryptor = New([]string{"test-key"})

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
		// NOTE: key is not validated against email format
		// {
		// 	name:        "Invalid email",
		// 	data:        "Sensitive campaign data",
		// 	email:       "invalid-email",
		// 	expectError: true,
		// },
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
			encryptedData, err := testEncryptor.Encrypt(Data{Data: tt.data, Key: tt.email})
			if (err != nil) != tt.expectError {
				t.Fatalf("Encrypt() error = %v, expectError = %v", err, tt.expectError)
			}

			// Skip decryption test if encryption fails
			if tt.expectError {
				return
			}

			// Decrypt the data
			decryptedData, err := testEncryptor.Decrypt(Data{Data: encryptedData, Key: tt.email})
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
	emptyEncryptor := New([]string{})

	data := "Sensitive campaign data"
	email := "test@example.com"

	_, err := emptyEncryptor.Encrypt(Data{Data: data, Key: email})
	if err == nil {
		t.Fatalf("Expected error for missing secret key, but got none")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	originalEncryptor := testEncryptor
	wrongKeyEncryptor := New([]string{"wrong-secure-secret-key"})

	data := "Sensitive campaign data"
	email := "test@example.com"

	// Encrypt with the original key
	encryptedData, err := originalEncryptor.Encrypt(Data{Data: data, Key: email})
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Try to decrypt with the wrong key
	_, err = wrongKeyEncryptor.Decrypt(Data{Data: encryptedData, Key: email})
	if err == nil {
		t.Fatalf("Expected error for decrypting with the wrong key, but got none")
	}
}
