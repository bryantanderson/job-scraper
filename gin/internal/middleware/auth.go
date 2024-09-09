package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	var isDevEnvironment bool
	environment := os.Getenv("ENVIRONMENT")
	apiToken := os.Getenv("API_TOKEN")

	if environment == "dev" {
		isDevEnvironment = true
	}

	return func(c *gin.Context) {
		// If authentication is disabled for development, skip
		if isDevEnvironment {
			c.Next()
			return
		}
		// Get authorization header from HTTP headers
		authorizationHeader := c.GetHeader("Authorization")
		token := strings.Split(authorizationHeader, " ")[1]

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		} else if token != apiToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}
		c.Next()
	}
}
