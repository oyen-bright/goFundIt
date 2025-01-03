package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
)

func CampaignKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		campaignKey := c.GetHeader("Campaign-Key")

		if campaignKey == "" {
			handlers.Unauthorized(c, "Unauthorized", nil)
			c.Abort()
			return
		}
		c.Set("Campaign-Key", campaignKey)

		c.Next()
	}
}
