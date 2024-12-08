package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		c.Next()
	}
}
