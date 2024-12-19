package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

// Helper functions

// getClaimsFromContext extracts JWT claims from the context
func getClaimsFromContext(c *gin.Context) jwt.Claims {
	return c.MustGet("claims").(jwt.Claims)
}

// GetCampaignID extracts the campaignId form the requestParam
func GetCampaignID(c *gin.Context) string {
	return c.Param("campaignID")
}

// parseActivityID converts the activity ID from the URL parameter to uint
func parseActivityID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("activityID"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// parseContributorID converts the contributor ID from the URL parameter to uint
func parseContributorID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("contributorID"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
