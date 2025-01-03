package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/api/handlers"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
)

func Auth(jwt jwt.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			handlers.Unauthorized(c, "Unauthorized", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			handlers.Unauthorized(c, "Unauthorized", nil)
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			handlers.Unauthorized(c, "Unauthorized", nil)
			c.Abort()
			return
		}
		c.Set("claims", *claims)
		c.Next()

	}

}
