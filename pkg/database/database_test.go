package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Skip if not in CI/Test environment
	if os.Getenv("TEST_DB") != "true" {
		t.Skip("Skipping database tests. Set TEST_DB=true to run")
	}

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: Config{
				Host:     "localhost",
				User:     "test_user",
				Password: "test_password",
				DBName:   "test_db",
				Port:     5432,
			},
			wantErr: false,
		},
		{
			name: "invalid configuration",
			config: Config{
				Host:     "invalid_host",
				User:     "test_user",
				Password: "test_password",
				DBName:   "test_db",
				Port:     5432,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Init(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, db)

			// Test closing the connection
			err = Close(db)
			assert.NoError(t, err)
		})
	}
}
