package handlers

import (
	"net/http"
	"unicode"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func handleInternalError(c *gin.Context, err error) {
	log.Errorln(err)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func handleBadRequest(c *gin.Context, err error) {
	log.Errorln(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Input")
}

func makeData(d any) gin.H {
	return gin.H{"data": d}
}

func firstToLower(s string) string {
	// Decode first letter to a rune and get size (in bytes)
	r, size := utf8.DecodeRuneInString(s)
	// Check that the decoding from string to runes was successful
	if r == utf8.RuneError && size < 1 {
		return s
	}

	lc := unicode.ToLower(r)

	// If the string is already lowercase, return
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}
