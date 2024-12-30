package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestContributorRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	contributor := &models.Contributor{
		Name:       "Test User",
		Email:      "test@example.com",
		CampaignID: "test-campaign",
		Amount:     100,
	}

	err := repo.Create(contributor)
	assert.NoError(t, err)
	assert.NotZero(t, contributor.ID)
}

func TestContributorRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	// Create initial contributor
	contributor := &models.Contributor{
		Name:       "Initial Name",
		Email:      "test@example.com",
		CampaignID: "test-campaign",
		Amount:     100,
	}
	err := repo.Create(contributor)
	assert.NoError(t, err)

	// Update contributor
	contributor.Name = "Updated Name"
	contributor.Amount = 200
	err = repo.Update(contributor)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetContributorById(contributor.ID, false)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, float64(200), updated.Amount)
}

func TestContributorRepository_GetContributorsByCampaignID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	// Create test contributors
	contributors := []models.Contributor{
		{Name: "User1", CampaignID: "campaign1", Amount: 100, Email: "contributor1@example.com"},
		{Name: "User2", CampaignID: "campaign1", Amount: 200, Email: "contributor2@example.com"},
		{Name: "User3", CampaignID: "campaign2", Amount: 300, Email: "contributor3@example.com"},
	}

	for _, c := range contributors {
		err := repo.Create(&c)
		assert.NoError(t, err)
	}

	// Test retrieval
	result, err := repo.GetContributorsByCampaignID("campaign1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestContributorRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	contributor := &models.Contributor{
		Name:       "Test User",
		Amount:     100,
		CampaignID: "test-campaign",
	}

	err := repo.Create(contributor)
	assert.NoError(t, err)

	err = repo.Delete(contributor)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetContributorById(contributor.ID, false)
	assert.Error(t, err)
}

func TestContributorRepository_GetEmailsOfActiveContributors(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	// Create test contributors
	contributors := []models.Contributor{
		{Name: "User1", CampaignID: "campaign1", Amount: 100, Email: "user1@example.com"},
		{Name: "User2", CampaignID: "campaign1", Amount: 200, Email: "user2@example.com"},
		{Name: "User3", CampaignID: "campaign2", Amount: 300, Email: "contributor3@example.com"},
	}

	for _, c := range contributors {
		err := repo.Create(&c)
		assert.NoError(t, err)
	}

	emails := []string{"user1@example.com", "user2@example.com", "nonexistent@example.com"}
	existingEmails, err := repo.GetEmailsOfActiveContributors(emails)

	assert.NoError(t, err)
	assert.Len(t, existingEmails, 2)
	assert.Contains(t, existingEmails, "user1@example.com")
	assert.Contains(t, existingEmails, "user2@example.com")
}

func TestContributorRepository_UpdateName(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewContributorRepository(db)

	contributor := &models.Contributor{
		Name: "Initial Name",

		CampaignID: "campaign1",
		Amount:     100,
		Email:      "user1@example.com",
	}
	err := repo.Create(contributor)
	assert.NoError(t, err)
	t.Logf("%+v", contributor)

	err = repo.UpdateName(contributor.ID, "New Name")
	assert.NoError(t, err)

	updated, err := repo.GetContributorById(contributor.ID, false)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
}
