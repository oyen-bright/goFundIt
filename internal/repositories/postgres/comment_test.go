package postgress

import (
	"testing"

	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCommentRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	comment := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)

	err = repo.Create(comment)
	assert.NoError(t, err)
	assert.NotEmpty(t, comment.ID)
}

func TestCommentRepository_Get(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	comment := &models.Comment{
		Content:    "Test comment",
		CreatedBy:  *user,
		ActivityID: 1,
	}

	err = db.Create(comment).Error
	assert.NoError(t, err)

	fetchedComment, err := repo.Get(comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content, fetchedComment.Content)
	assert.Equal(t, comment.ID, fetchedComment.ID)
}

func TestCommentRepository_GetByActivityID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	comment1 := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)
	comment2 := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)
	err = db.Create(comment1).Error
	assert.NoError(t, err)
	err = db.Create(comment2).Error
	assert.NoError(t, err)

	comments, err := repo.GetByActivityID(1)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
}

func TestCommentRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	comment := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)

	err = db.Create(comment).Error
	assert.NoError(t, err)

	comment.Content = "Updated content"
	err = repo.Update(comment)
	assert.NoError(t, err)

	updatedComment, err := repo.Get(comment.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Updated content", updatedComment.Content)
}

func TestCommentRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	comment := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)

	err = db.Create(comment).Error
	assert.NoError(t, err)

	err = repo.Delete(comment.ID)
	assert.NoError(t, err)

	var deletedComment models.Comment
	err = db.First(&deletedComment, comment.ID).Error
	assert.Error(t, err)
}

func TestCommentRepository_FindReplies(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewCommentRepository(db)

	user, err := createTestUser(db)
	assert.NoError(t, err)

	parentComment := models.NewComment(
		nil,
		1,
		"Test comment",
		*user,
	)

	err = db.Create(parentComment).Error
	assert.NoError(t, err)

	reply := models.NewComment(
		&parentComment.ID,
		1,
		"Test comment",
		*user,
	)

	err = db.Create(reply).Error
	assert.NoError(t, err)

	replies, err := repo.FindReplies(parentComment.ID)
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, reply.Content, replies[0].Content)
}
