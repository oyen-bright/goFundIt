package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CampaignKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignKey := c.GetHeader("Campaign-Key")

		if campaignKey == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
			return
		}
		c.Set("Campaign-Key", campaignKey)

		c.Next()
	}
}
