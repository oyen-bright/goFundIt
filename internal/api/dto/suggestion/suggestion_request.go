package dto

// SuggestionRequest represents the request body for getting suggestions
type SuggestionRequest struct {
	// Content for generating suggestions
	Content string `json:"content" binding:"required" example:"Trip to London for a day at midnight"`
}
