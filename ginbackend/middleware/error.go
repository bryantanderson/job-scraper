package middleware

import (
	"log"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, error := range c.Errors {
		log.Fatal(error)
	}
}
