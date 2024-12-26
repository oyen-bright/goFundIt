package dto

type ContributorEditRequest struct {
	Name string `json:"name" binding:"required,gte=3"`
}
