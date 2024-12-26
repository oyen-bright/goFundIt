package interfaces

import "github.com/oyen-bright/goFundIt/internal/models"

type CommentService interface {
	CreateComment(comment *models.Comment, campaignID string, activityID uint, userHandle string) error
	DeleteComment(commentID, userHandle string) error
	UpdateComment(comment models.Comment, userHandle string) error

	GetActivityComments(activityID uint) ([]models.Comment, error)
	GetCommentReplies(commentID string) ([]models.Comment, error)
}
