package postgress_test

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	postgress "github.com/oyen-bright/goFundIt/internal/repositories/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO:setup a postgress test database
// setupTestDB creates a new in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.Contributor{})
	require.NoError(t, err)

	return db
}

// createTestUser creates a test user with given parameters
func createTestUser(email, name, handle string) *models.User {
	return &models.User{
		Email:    email,
		Name:     name,
		Handle:   handle,
		Verified: false,
	}
}

func TestAuthRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := postgress.NewAuthRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		update  *models.User // for testing update scenario
		wantErr bool
	}{
		{
			name:    "successful save new user",
			user:    createTestUser("test@example.com", "Test User", "test-handle"),
			wantErr: false,
		},
		{
			name:    "update existing user",
			user:    createTestUser("update@example.com", "Original Name", "update-handle"),
			update:  createTestUser("update@example.com", "Updated Name", "update-handle"),
			wantErr: false,
		},
		{
			name:    "save user with empty email",
			user:    createTestUser("", "Test User", "test-handle"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initial save
			err := repo.Save(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Test update if provided
			if tt.update != nil {
				err = repo.Save(tt.update)
				assert.NoError(t, err)

				// Verify update
				var updatedUser models.User
				err = db.First(&updatedUser, "email = ?", tt.update.Email).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.update.Name, updatedUser.Name)
			}
		})
	}
}

func TestAuthRepository_CreateMultiple(t *testing.T) {
	db := setupTestDB(t)
	repo := postgress.NewAuthRepository(db)

	tests := []struct {
		name    string
		users   []models.User
		wantLen int
		wantErr bool
	}{
		{
			name: "successful batch create",
			users: []models.User{
				*createTestUser("user1@example.com", "User 1", "handle-1"),
				*createTestUser("user2@example.com", "User 2", "handle-2"),
				*createTestUser("user3@example.com", "User 3", "handle-3"),
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name:    "empty users list",
			users:   []models.User{},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "duplicate emails",
			users: []models.User{
				*createTestUser("duplicate@example.com", "User 1", "handle-1"),
				*createTestUser("duplicate@example.com", "User 2", "handle-2"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUsers, err := repo.CreateMultiple(tt.users)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, createdUsers, tt.wantLen)

			// Verify users were created in database
			if tt.wantLen > 0 {
				var count int64
				db.Model(&models.User{}).Count(&count)
				assert.Equal(t, int64(tt.wantLen), count)
			}
		})
	}
}

func TestAuthRepository_FindByHandle(t *testing.T) {
	db := setupTestDB(t)
	repo := postgress.NewAuthRepository(db)

	// Create test user with contributions
	testUser := createTestUser("test@example.com", "Test User", "test-handle")
	err := db.Create(testUser).Error
	require.NoError(t, err)

	// Create some contributions for the user
	contributions := []models.Contributor{
		{Email: testUser.Email, Amount: 100},
		{Email: testUser.Email, Amount: 200},
	}
	err = db.Create(&contributions).Error
	require.NoError(t, err)

	tests := []struct {
		name              string
		handle            string
		wantEmail         string
		wantContributions int
		wantErr           bool
	}{
		{
			name:              "existing user",
			handle:            "test-handle",
			wantEmail:         "test@example.com",
			wantContributions: 2,
			wantErr:           false,
		},
		{
			name:    "non-existing user",
			handle:  "non-existing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByHandle(tt.handle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tt.wantEmail, user.Email)
			assert.Len(t, user.Contributions, tt.wantContributions)
		})
	}
}

// func TestAuthRepository_FindNonExistingUsers(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := postgress.NewAuthRepository(db)

// 	// Create existing user
// 	existingUser := createTestUser("existing@example.com", "Existing User", "existing-handle")
// 	err := db.Create(existingUser).Error
// 	require.NoError(t, err)

// 	tests := []struct {
// 		name      string
// 		users     []models.User
// 		wantCount int
// 		wantErr   bool
// 	}{
// 		{
// 			name: "all non-existing users",
// 			users: []models.User{
// 				*createTestUser("new1@example.com", "New User 1", "new-handle-1"),
// 				*createTestUser("new2@example.com", "New User 2", "new-handle-2"),
// 			},
// 			wantCount: 2,
// 			wantErr:   false,
// 		},
// 		{
// 			name: "mix of existing and non-existing users",
// 			users: []models.User{
// 				*createTestUser("existing@example.com", "Existing User", "existing-handle"),
// 				*createTestUser("new3@example.com", "New User 3", "new-handle-3"),
// 			},
// 			wantCount: 1,
// 			wantErr:   false,
// 		},
// 		{
// 			name:      "empty user list",
// 			users:     []models.User{},
// 			wantCount: 0,
// 			wantErr:   false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			nonExistingUsers, err := repo.FindNonExistingUsers(tt.users)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Len(t, nonExistingUsers, tt.wantCount)

// 			// Verify returned users are actually non-existing
// 			for _, user := range nonExistingUsers {
// 				var count int64
// 				db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
// 				assert.Equal(t, int64(0), count)
// 			}
// 		})
// 	}
// }

func TestAuthRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := postgress.NewAuthRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name:    "successful delete",
			user:    createTestUser("delete@example.com", "Delete User", "delete-handle"),
			wantErr: false,
		},
		{
			name:    "delete non-existing user",
			user:    &models.User{}, // empty user
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First create the user if it's a valid test user
			if tt.user.Email != "" {
				err := db.Create(tt.user).Error
				require.NoError(t, err)
			}

			err := repo.Delete(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify user was deleted
			var count int64
			db.Model(&models.User{}).Where("email = ?", tt.user.Email).Count(&count)
			assert.Equal(t, int64(0), count)
		})
	}
}
