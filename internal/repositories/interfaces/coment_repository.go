package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type CommentRepository interface {
	Create(comment *models.Comment) error
	Get(commentID string) (models.Comment, error)
	GetByActivityID(activityID uint) ([]models.Comment, error)
	Delete(commentID string) error
	FindReplies(commentID string) ([]models.Comment, error)
	Update(comment *models.Comment) error
}
