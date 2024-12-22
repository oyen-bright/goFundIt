package interfaces

import (
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/payout"
	"github.com/oyen-bright/goFundIt/internal/models"
)

type PayoutService interface {
	InitializePayout(campaignID, userHandle string, req dto.PayoutRequest) (*models.Payout, error)
	//TODO:change response to DTO
	GetBankList() ([]interface{}, error)
	//TODO:change response to DTO
	VerifyAccount(dto.VerifyAccountRequest) (interface{}, error)
	GetPayoutByCampaignID(campaignID string) (*models.Payout, error)
}
