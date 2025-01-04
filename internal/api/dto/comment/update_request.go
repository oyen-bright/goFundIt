package dto

// UpdateCommentRequest represents the request body for creating a comment
// @Description Request structure for creating a new comment
type UpdateCommentRequest struct {
	// Content of the comment
	// @example "This is a great project! Looking forward to contributing."
	Content string `json:"content" binding:"required"`
}
