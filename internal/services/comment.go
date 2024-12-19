package services

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/database"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type commentService struct {
	repo            interfaces.CommentRepository
	authService     services.AuthService
	activityService services.ActivityService
	logger          logger.Logger
}

func NewCommentService(repo interfaces.CommentRepository, authService services.AuthService, activityService services.ActivityService, logger logger.Logger) services.CommentService {
	return &commentService{
		repo:            repo,
		authService:     authService,
		activityService: activityService,
		logger:          logger,
	}
}

// CreateComment creates a new comment on a given activityID
func (c *commentService) CreateComment(comment *models.Comment, campaignID string, activityID uint, userHandle string) error {

	// Validate user
	user, err := c.authService.GetUserByHandle(userHandle)
	if err != nil {
		return err
	}

	// Validate activity
	if _, err := c.activityService.GetActivityByID(activityID, campaignID); err != nil {
		return err
	}

	if comment.ParentID != nil && *comment.ParentID != "" {
		// Validate parent comment
		_, err := c.repo.Get(*comment.ParentID)
		if err != nil {
			if database.Error(err).IsNotfound() {
				return errs.BadRequest("Parent comment not found", err)
			}
			return errs.InternalServerError(err).Log(c.logger)
		}
	}

	comment.FromBinding(user, activityID)
	err = c.repo.Create(comment)

	if err != nil {
		if database.Error(err).IsPrimaryKeyViolated() {
			return errs.BadRequest("parent comment not found", err)
		}

		return errs.InternalServerError(err).Log(c.logger)
	}
	return nil
}

// DeleteComment deletes a comment by ID
func (c *commentService) DeleteComment(commentID string, userHandle string) error {
	//Validate comment for modification
	_, err := c.validateCommentForModification(commentID, userHandle)
	if err != nil {
		return err
	}
	//Delete comment
	return c.repo.Delete(commentID)
}

func (c *commentService) GetActivityComments(activityID uint) ([]models.Comment, error) {
	comments, err := c.repo.GetByActivityID(activityID)
	if err != nil {
		return nil, errs.InternalServerError(err).Log(c.logger)
	}
	return comments, nil
}

// GetCommentReplies gets all replies to a comment
func (c *commentService) GetCommentReplies(commentID string) ([]models.Comment, error) {
	comments, err := c.repo.FindReplies(commentID)
	if err != nil {
		return comments, errs.InternalServerError(err).Log(c.logger)
	}
	return comments, nil
}

// UpdateComment updates a comment by the comment ID
func (c *commentService) UpdateComment(comment models.Comment, userHandle string) error {
	//Validate comment for modification
	_, err := c.validateCommentForModification(comment.ID, userHandle)
	if err != nil {
		return err
	}

	// Update comment
	return c.repo.Update(&comment)

}

// Helper function

func (c *commentService) validateCommentForModification(commentID, userHandle string) (*models.Comment, error) {
	//Validate user
	user, err := c.authService.GetUserByHandle(userHandle)
	if err != nil {
		return nil, err
	}

	//Validate comment
	comment, err := c.repo.Get(commentID)
	if err != nil {
		if database.Error(err).IsNotfound() {
			return nil, errs.BadRequest("Comment not found", err)
		}
		return nil, errs.InternalServerError(err).Log(c.logger)
	}

	if comment.CreatedByHandle != user.Handle {
		return nil, errs.BadRequest("You can't modify a comment you didn't create", nil)
	}

	return &comment, nil
}
