package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
)

type CommentHandler struct {
	CommentService interfaces.CommentService
}

func NewCommentHandler(CommentService interfaces.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: CommentService,
	}

}

// @Summary Create Comment
// @Description Posts a new comment or reply on an activity
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param request body dto.CreateCommentRequest true "Comment Details"
// @Success 200 {object} SuccessResponse{data=models.Comment} "Comment created successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/comments [post]
// CreateComment creates a new comment on a given activityID
func (h *CommentHandler) HandleCreateComment(c *gin.Context) {
	var comment models.Comment
	claims := getClaimsFromContext(c)
	campaignID := GetCampaignID(c)
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	if err := c.BindJSON(&comment); err != nil {
		BadRequest(c, "Invalid inputs", ExtractValidationErrors(err))
		return
	}

	err = h.CommentService.CreateComment(&comment, campaignID, activityID, claims.Handle)

	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Comment create successfully", comment)
}

// @Summary Delete Comment
// @Description Deletes a specific comment
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param commentID path string true "Comment ID"
// @Success 200 {object} SuccessResponse "Comment deleted successfully"
// @Failure 400 {object} BadRequestResponse "Invalid request"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Comment not found"
// @Router /activity/{campaignID}/{activityID}/comments/{commentID} [delete]
// HandleDeleteComment handles the deletion of a comment
func (h *CommentHandler) HandleDeleteComment(c *gin.Context) {
	commentID := getCommentID(c)
	claims := getClaimsFromContext(c)

	err := h.CommentService.DeleteComment(commentID, claims.Handle)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Comment deleted successfully", nil)
}

// @Summary Get Activity Comments
// @Description Retrieves all comments associated with a specific activity
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Success 200 {object} SuccessResponse{data=[]models.Comment} "Comments retrieved successfully"
// @Failure 400 {object} BadRequestResponse "Invalid Activity ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Router /activity/{campaignID}/{activityID}/comments [get]
// HandleGetActivityComments handles the retrieval of comments for a given activity
func (h *CommentHandler) HandleGetActivityComments(c *gin.Context) {
	activityID, err := parseActivityID(c)
	if err != nil {
		BadRequest(c, "Invalid Activity ID", nil)
		return
	}

	comments, err := h.CommentService.GetActivityComments(activityID)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Comments retrieved successfully", comments)
}

// @Summary Get replies for a comment
// @Description Retrieves all reply comments for a given parent comment ID
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} SuccessResponse{data=[]models.Comment} "Successfully retrieved replies"
// @Failure 400 {object} BadRequestResponse "Invalid comment ID"
// @Failure 404 {object} BadRequestResponse "Comment not found"
// @Failure 500 {object} response "Internal server error"
// @Router /comments/{comment_id}/replies [get]
// HandleGetCommentReplies handles the retrieval of replies to a comment
func (h *CommentHandler) HandleGetCommentReplies(c *gin.Context) {
	commentID := getCommentID(c)
	comments, err := h.CommentService.GetCommentReplies(commentID)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Replies retrieved successfully", comments)
}

// @Summary Update Comment
// @Description Modifies the content of an existing comment
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Param activityID path string true "Activity ID"
// @Param commentID path string true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Updated Comment Content"
// @Success 200 {object} SuccessResponse{data=models.Comment} "Comment updated successfully"
// @Failure 400 {object} BadRequestResponse{errors=[]ValidationError} "Invalid inputs"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Comment not found"
// @Router /activity/{campaignID}/{activityID}/comments/{commentID} [patch]
// HandleUpdateComment handles the update of a comment
func (h *CommentHandler) HandleUpdateComment(c *gin.Context) {
	var comment models.Comment
	claims := getClaimsFromContext(c)
	commentID := getCommentID(c)

	if err := c.BindJSON(&comment); err != nil {
		BadRequest(c, "Invalid inputs", ExtractValidationErrors(err))
		return
	}
	comment.ID = commentID

	err := h.CommentService.UpdateComment(comment, claims.Handle)
	if err != nil {
		FromError(c, err)
		return
	}
	Success(c, "Comment updated successfully", comment)
}
