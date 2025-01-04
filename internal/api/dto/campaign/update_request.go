package dto

import "time"

type CampaignUpdateRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
}
