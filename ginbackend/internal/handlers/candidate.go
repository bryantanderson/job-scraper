package handlers

import (
	"net/http"

	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) registerCandidateRoutes() {
	candidateGroup := s.Router.Group("/candidates")
	{
		candidateGroup.POST("/", s.handleCandidateCreate())
		candidateGroup.GET("/:id", s.handleCandidateGet())
		candidateGroup.DELETE("/:id", s.handleCandidateDelete())
	}
}

// handleCandidateCreate godoc
// @Summary Creates a candidate.
// @Tags Candidate
// @Description Creates a candidate.
// @Param data body services.CandidateDto true "Request payload"
// @Accept json
// @Produce json
// @Success 201 {object} services.Candidate
// @Failure 400
// @Failure 500
// @Router /candidates/ [post]
func (s *Server) handleCandidateCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto services.CandidateDto

		if err := c.ShouldBindJSON(&dto); err != nil {
			handleBadRequest(c, err)
			return
		}

		candidate, err := s.CandidateService.CreateCandidate(&dto)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, candidate)
	}
}

// handleCandidateGet godoc
// @Summary Gets a candidate.
// @Tags Candidate
// @Description Gets a candidate.
// @Param id path string true "Candidate ID"
// @Accept json
// @Produce json
// @Success 201 {object} services.Candidate
// @Failure 404
// @Failure 500
// @Router /candidates/{id} [get]
func (s *Server) handleCandidateGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		candidateId := c.Param("id")
		candidate, err := s.CandidateService.GetCandidate(candidateId)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "Candidate not found")
			return
		}

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, candidate)
	}
}

// handleCandidateDelete godoc
// @Summary Deletes an existing candidate.
// @Tags Candidate
// @Description Deletes an existing candidate.
// @Param id path string true "Candidate ID"
// @Accept json
// @Success 204
// @Failure 500
// @Router /candidates/{id} [delete]
func (s *Server) handleCandidateDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		candidateId := c.Param("id")
		err := s.CandidateService.DeleteCandidate(candidateId)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusNoContent, "Candidate successfully deleted")
	}
}
