package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/response"
	"github.com/oyen-bright/goFundIt/pkg/utils"
)

type CommentHandler struct {
	CommentService interfaces.CommentService
}

func NewCommentHandler(CommentService interfaces.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: CommentService,
	}

}

// CreateComment creates a new comment on a given activityID
func (h *CommentHandler) HandleCreateComment(c *gin.Context) {
	var comment models.Comment
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	activityID, err := parseActivityID(c)
	if err != nil {
		response.BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := c.BindJSON(&comment); err != nil {
		response.BadRequest(c, "Invalid inputs", utils.ExtractValidationErrors(err))
		return
	}

	err = h.CommentService.CreateComment(&comment, campaignID, activityID, claims.Handle)

	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Comment create successfully", comment)
}

// HandleDeleteComment handles the deletion of a comment
func (h *CommentHandler) HandleDeleteComment(c *gin.Context) {
	commentID := getCommentID(c)
	claims := getClaimsFromContext(c)

	err := h.CommentService.DeleteComment(commentID, claims.Handle)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Comment deleted successfully", nil)
}

// HandleGetActivityComments handles the retrieval of comments for a given activity
func (h *CommentHandler) HandleGetActivityComments(c *gin.Context) {
	activityID, err := parseActivityID(c)
	if err != nil {
		response.BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	comments, err := h.CommentService.GetActivityComments(activityID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Comments retrieved successfully", comments)
}

// HandleGetCommentReplies handles the retrieval of replies to a comment
func (h *CommentHandler) HandleGetCommentReplies(c *gin.Context) {
	commentID := getCommentID(c)
	comments, err := h.CommentService.GetCommentReplies(commentID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Replies retrieved successfully", comments)
}

// HandleUpdateComment handles the update of a comment
func (h *CommentHandler) HandleUpdateComment(c *gin.Context) {
	var comment models.Comment
	claims := getClaimsFromContext(c)
	commentID := getCommentID(c)

	if err := c.BindJSON(&comment); err != nil {
		response.BadRequest(c, "Invalid inputs", utils.ExtractValidationErrors(err))
		return
	}
	comment.ID = commentID

	err := h.CommentService.UpdateComment(comment, claims.Handle)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, "Comment updated successfully", comment)
}
