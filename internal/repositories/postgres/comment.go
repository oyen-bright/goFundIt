package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) interfaces.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (c *commentRepository) Create(comment *models.Comment) error {
	return c.db.Create(comment).Error
}
func (c *commentRepository) Get(commentID string) (models.Comment, error) {

	var comment models.Comment

	err := c.db.Where("id = ?", commentID).Preload("CreatedBy").First(&comment).Error
	return comment, err
}

func (r *commentRepository) GetByActivityID(activityID uint) ([]models.Comment, error) {
	var comments []models.Comment

	err := r.db.
		Preload("CreatedBy").
		Preload("Replies.CreatedBy").
		Where("activity_id = ? AND parent_id IS NULL", activityID).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, err
}

func (c *commentRepository) Delete(commentID string) error {
	return c.db.Where("id = ?", commentID).Delete(&models.Comment{}).Error
}

func (c *commentRepository) FindReplies(commentID string) ([]models.Comment, error) {
	var comments []models.Comment

	err := c.db.Where("parent_id = ?", commentID).Preload("CreatedBy").Find(&comments).Error
	return comments, err
}

func (c *commentRepository) Update(comment *models.Comment) error {
	return c.db.Model(&models.Comment{}).Where("id = ?", comment.ID).Updates(
		map[string]interface{}{
			"content": comment.Content,
		}).Error
}
