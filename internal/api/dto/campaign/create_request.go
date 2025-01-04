package dto

import (
	"time"

	"github.com/oyen-bright/goFundIt/internal/models"
)

// CampaignRequest represents the campaign creation/update payload
// @Description Campaign creation request structure
type CampaignRequest struct {
	// @Description Campaign title
	// @example "Save the Forests"
	Title string `json:"title" binding:"required" validate:"required,min=4"`

	// @Description Detailed campaign description
	// @example "This campaign aims to protect and preserve our forest ecosystems..."
	Description string `json:"description" binding:"required,min=100" validate:"required,min=100"`

	// @Description Payment method (fiat/crypto/manual)
	// @example "fiat"
	PaymentMethod models.PaymentMethod `json:"paymentMethod" binding:"required,oneof=fiat crypto manual" validate:"required,oneof=fiat crypto manual"`

	// @Description Fiat currency code (required if paymentMethod is fiat)
	// @example "USD"
	FiatCurrency *models.FiatCurrency `json:"fiatCurrency,omitempty" binding:"required_if=PaymentMethod fiat" validate:"required_if=PaymentMethod fiat,omitempty"`

	// @Description Crypto token (required if paymentMethod is crypto)
	// @example "ETH"
	CryptoToken *models.CryptoToken `json:"cryptoToken,omitempty" binding:"required_if=PaymentMethod crypto" validate:"required_if=PaymentMethod crypto,omitempty"`

	// @Description Campaign images
	Images []models.CampaignImage `json:"images" binding:"omitempty,dive,required"`

	// @Description Campaign contributors
	Contributors []models.Contributor `json:"contributors" binding:"required,gt=0,dive,required" validate:"required,gt=0,dive,required"`

	// @Description CampaignActivities
	Activities []models.Activity `json:"activities" binding:"omitempty,dive,required" validate:"omitempty,dive,required"`

	// @Description Campaign start date
	// @example "2025-01-03T13:51:06Z"
	StartDate time.Time `json:"startDate" binding:"required" validate:"required"`

	// @Description Campaign end date (must be after start date)
	// @example "2025-02-03T13:51:06Z"
	EndDate time.Time `json:"endDate" binding:"required,gtfield=StartDate" validate:"required,gtfield=StartDate"`
}
