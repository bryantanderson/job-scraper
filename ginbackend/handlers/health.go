package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) registerHealthRoutes() {
	s.Router.GET("/health", s.handleHealthCheck())
}

// Health Check godoc
// @Summary Health check
// @Tags health
// @Description Health check ping
// @Accept json
// @Produce json
// @Success 200
// @Router /health [get]
func (s *Server) handleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Health OK"})
	}
}
