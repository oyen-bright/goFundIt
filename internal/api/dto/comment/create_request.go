package dto

// CreateCommentRequest represents the request body for creating a comment
// @Description Request structure for creating a new comment
type CreateCommentRequest struct {
	// Content of the comment
	// @example "This is a great project! Looking forward to contributing."
	Content string `json:"content" binding:"required"`

	// ID of the activity this comment belongs to
	// @example 123
	ActivityID uint `json:"activityID" binding:"required"`

	// Optional ID of the parent comment if this is a reply
	// @example "comment-123-abc"
	ParentID *string `json:"parentId,omitempty"`
}
