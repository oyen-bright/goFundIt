package middlewares

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PaystackSignature(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read request body
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Get signature from header
		signature := c.GetHeader("x-paystack-signature")
		if signature == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Compute HMAC hash
		mac := hmac.New(sha512.New, []byte(secretKey))
		mac.Write(payload)
		hash := hex.EncodeToString(mac.Sum(nil))

		// Verify signature
		if hash != signature {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
