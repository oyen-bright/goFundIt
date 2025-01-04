package dto

type AuthRequest struct {
	// @Description User's email
	// @example
	Email string `json:"email" binding:"required,email"`

	// @Description Name of the user
	// @example
	Name string `json:"name" binding:"omitempty"`
}
