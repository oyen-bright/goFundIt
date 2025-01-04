package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
)

func createTestActivity() *models.Activity {
	return &models.Activity{
		Title:           "Test Activity",
		Subtitle:        "Test Subtitle",
		ImageUrl:        "https://test.com/image.jpg",
		IsMandatory:     true,
		Cost:            100.0,
		IsApproved:      true,
		CreatedByHandle: "test_handle",
		CampaignID:      "campaign123",
	}
}

func TestActivityRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	created, err := repo.Create(activity)

	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, activity.Title, created.Title)
}

func TestActivityRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	created, err := repo.Create(activity)
	assert.NoError(t, err)

	found, err := repo.GetByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Title, found.Title)
}

func TestActivityRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	created, err := repo.Create(activity)
	assert.NoError(t, err)

	created.Title = "Updated Title"
	err = repo.Update(&created)
	assert.NoError(t, err)

	updated, err := repo.GetByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
}

func TestActivityRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	created, err := repo.Create(activity)
	assert.NoError(t, err)

	err = repo.Delete(&created)
	assert.NoError(t, err)

	_, err = repo.GetByID(created.ID)
	assert.Error(t, err)
}

func TestActivityRepository_GetByCampaignID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	_, err := repo.Create(activity)
	assert.NoError(t, err)

	activities, err := repo.GetByCampaignID(activity.CampaignID)
	assert.NoError(t, err)
	assert.Len(t, activities, 1)
	assert.Equal(t, activity.Title, activities[0].Title)
}

func TestActivityRepository_AddAndRemoveContributor(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewActivityRepo(db)

	activity := createTestActivity()
	created, err := repo.Create(activity)
	assert.NoError(t, err)

	contributor := &models.Contributor{Name: "Test Contributor", CampaignID: created.CampaignID, Amount: 100}
	err = db.Create(contributor).Error
	assert.NoError(t, err)

	err = repo.AddContributor(created.ID, contributor.ID)
	assert.NoError(t, err)

	participants, err := repo.GetParticipants(created.ID)
	assert.NoError(t, err)
	assert.Len(t, participants, 1)
	assert.Equal(t, contributor.Name, participants[0].Name)

	err = repo.RemoveContributor(created.ID, contributor.ID)
	assert.NoError(t, err)

	participants, err = repo.GetParticipants(created.ID)
	assert.NoError(t, err)
	assert.Len(t, participants, 0)
}
