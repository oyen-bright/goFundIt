package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthSave(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewAuthRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: models.NewUser("test", "test@example.com", true),

			wantErr: false,
		},
		{
			name:    "duplicate email",
			user:    models.NewUser("test", "test@example.com", true),
			wantErr: false, // Should not error due to OnConflict clause
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthFindByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewAuthRepository(db)

	testUser := models.NewUser("test", "test@example.com", true)

	err := repo.Save(testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			email:   "nonexistent@example.com",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestAuthFindExistingAndNonExistingUsers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewAuthRepository(db)

	// Create test users
	testUsers := []models.User{
		*models.NewUser("user1", "user1@example.com", true),
		*models.NewUser("user2", "user2@example.com", true),
	}
	users, err := repo.CreateMultiple(testUsers)
	println(users)
	require.NoError(t, err)

	emails := []string{
		"user1@example.com",
		"user2@example.com",
		"nonexistent@example.com",
	}

	existing, nonExisting, err := repo.FindExistingAndNonExistingUsers(emails)
	assert.NoError(t, err)
	assert.Len(t, existing, 2)
	assert.Len(t, nonExisting, 1)
	assert.Contains(t, nonExisting, "nonexistent@example.com")
}

func TestAuthCreateMultiple(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewAuthRepository(db)

	users := []models.User{
		*models.NewUser("user1", "user1@example.com", true),
		*models.NewUser("user2", "user2@example.com", true),
	}

	created, err := repo.CreateMultiple(users)
	assert.NoError(t, err)
	assert.Len(t, created, 2)

	// Verify users were created
	for _, user := range users {
		found, err := repo.FindByEmail(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.Name, found.Name)
	}
}
