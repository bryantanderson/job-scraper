package handlers

import (
	"net/http"
	"reflect"
	"sincidium/linkd/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) registerAssessRoutes() {
	assessGroup := s.Router.Group("/assessments")
	{
		assessGroup.GET("/", s.handleAssessmentQuery())
		assessGroup.GET("/:userId", s.handleAssessmentGet())
	}
}

// handleAssessmentGet godoc
// @Summary Gets an existing assessment.
// @Tags Assessment
// @Description Gets an existing assessment.
// @Param userId path string true "User ID"
// @Accept json
// @Produce json
// @Success 200 {object} services.Assessment
// @Failure 404
// @Failure 500
// @Router /assessments/{userId} [get]
func (s *Server) handleAssessmentGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		assessment, err := s.AssessorService.GetAssessment(userId)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "Assessment not found")
			return
		}

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, assessment)
	}
}

// handleAssessmentQuery godoc
// @Summary Queries for an existing assessment.
// @Tags Assessment
// @Description Queries for an existing assessment.
// @Param jobId query string false "ID of the job the assessment is based on"
// @Param score query string false "The score attained in the assessment"
// @Accept json
// @Produce json
// @Success 200 {array} services.Assessment
// @Failure 500
// @Router /assessments/ [get]
func (s *Server) handleAssessmentQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]string)

		// Get the type of the assessment, and the fields
		assessment := services.Assessment{}
		assessmentType := reflect.TypeOf(assessment)
		fields := reflect.VisibleFields(assessmentType)

		// Iterate through the fields and get non-empty query parameters
		for _, field := range fields {
			// Convert exported field name to lower
			key := firstToLower(field.Name)
			v, ok := c.GetQuery(key)

			if !ok || v == "" {
				continue
			}

			params[key] = v
		}

		assessments, err := s.AssessorService.QueryAssessments(params)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, assessments)
	}
}
