package campaign

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func CheckCampaignKey() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		campaignKey := c.GetHeader("Campaign-Key")

// 		if campaignKey == "" {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
// 			return
// 		}

// 		c.Next()
// 	}
// }
