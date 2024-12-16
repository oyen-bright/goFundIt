package campaign

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func validateContributionSum(fl validator.FieldLevel) bool {
	campaign, ok := fl.Parent().Interface().(Campaign)
	if !ok {
		return false
	}

	totalContributions := calculateContributionTotal(campaign)

	return totalContributions == campaign.TargetAmount
}

func calculateContributionTotal(campaign Campaign) float64 {
	var totalAmount float64
	for _, contributor := range campaign.Contributors {
		totalAmount += contributor.Amount
	}

	return totalAmount
}

func isCampaignStartDateValid(campaign Campaign) bool {
	return campaign.StartDate.After(time.Now()) || campaign.StartDate.Equal(time.Now())
}
