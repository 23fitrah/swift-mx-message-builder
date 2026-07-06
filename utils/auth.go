package utils

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewareGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusAccepted, gin.H{
				"response_code": "0002",
				"message":       "Missing or invalid Authorization header",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		authToken := os.Getenv("TOKEN_SWIFT_MX")

		if token != authToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"response_code": "0002",
				"message":       "Authentication failed",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
