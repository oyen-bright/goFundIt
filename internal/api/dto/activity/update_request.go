package dto

// UpdateActivityRequest represents the activity update payload
// @Description Activity update request structure
type UpdateActivityRequest struct {
	// Activity title (minimum 4 characters)
	// @example "Plant Trees in Local Park"
	Title string `json:"title" binding:"required,min=4"`

	// Optional subtitle for additional context
	// @example "Phase 1 of reforestation project"
	Subtitle string `json:"subtitle"`

	// Optional URL for activity image
	// @example "https://example.com/images/tree-planting.jpg"
	ImageUrl string `json:"image_url" binding:"omitempty,url" validate:"omitempty,url"`

	// Whether this activity is mandatory for the campaign
	// @example true
	IsMandatory bool `json:"is_mandatory" binding:"boolean"`

	// Cost of the activity (must be greater than 0)
	// @example 1500.50
	Cost float64 `json:"cost" binding:"required" validate:"required,gt=0"`

	// Approval status of the activity
	// @example false
	IsApproved bool `json:"is_approved"`
}
