package migrations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	if os.Getenv("TEST_DB") != "true" {
		t.Skip("Skipping database tests. Set TEST_DB=true to run")
	}

	dsn := "host=localhost user=test_user password=test_password dbname=test_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	sqlDB, err := db.DB()
	require.NoError(t, err)
	cleanUp := func() {
		sqlDB.Close()
	}

	return db, cleanUp
}

func TestMigrate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name    string
		db      *gorm.DB
		wantErr bool
	}{
		{
			name:    "successful migration",
			db:      db,
			wantErr: false,
		},
		{
			name:    "nil database",
			db:      nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Migrate(tt.db)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify tables exist
			var tables []string
			err = tt.db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Pluck("table_name", &tables).Error
			assert.NoError(t, err)
			assert.NotEmpty(t, tables)
		})
	}
}
