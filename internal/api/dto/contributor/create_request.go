package dto

// CreateContributorRequest represents the request body for creating a contributor
// @Description Request structure for creating a new contributor
type CreateContributorRequest struct {
	// Optional name of the contributor
	// @example "John Doe"
	Name string `json:"name" binding:"omitempty,gte=3"`

	// ID of the campaign the contribution is for
	// @example "campaign-123-abc"
	CampaignID string `json:"campaignId" binding:"required"`

	// Amount of the contribution
	// @example 100.50
	// @minimum 0
	Amount float64 `json:"amount" binding:"required,gte=0"`

	// Email address of the contributor
	// @example "john.doe@example.com"
	Email string `json:"email" binding:"required,email"`
}
