package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
)

func APIKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != key {
			handlers.Unauthorized(c, "Unauthorized", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
