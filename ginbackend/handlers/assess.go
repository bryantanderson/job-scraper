package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"reflect"
	"sincidium/linkd/api/services"
)

func (s *Server) registerAssessRoutes() {
	assessGroup := s.Router.Group("/assessments")
	{
		assessGroup.GET("/", s.handleAssessmentQuery())
		assessGroup.GET("/:userId", s.handleAssessmentGet())
	}
}

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

		c.JSON(http.StatusOK, makeData(assessment))
	}
}

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

		c.JSON(http.StatusOK, makeData(assessments))
	}
}
