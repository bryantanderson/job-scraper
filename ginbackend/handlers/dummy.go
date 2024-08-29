package handlers

import (
	"net/http"
	"sincidium/linkd/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) registerDummyRoutes() {
	dummyGroup := s.Router.Group("/dummies")
	{
		dummyGroup.POST("/", s.handleDummyCreate())
		dummyGroup.GET("/:id", s.handleDummyGet())
		dummyGroup.PATCH("/:id", s.handleDummyUpdate())
		dummyGroup.DELETE("/:id", s.handleDummyDelete())
	}
}

// handleDummyCreate godoc
// @Summary Creates a dummy.
// @Tags Dummy
// @Description Creates a dummy.
// @Param data body services.DummyDto true "Request payload"
// @Accept json
// @Produce json
// @Success 201 {object} services.Dummy
// @Failure 400
// @Failure 500
// @Router /dummies/ [post]
func (s *Server) handleDummyCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto services.DummyDto

		if err := c.ShouldBindJSON(&dto); err != nil {
			handleBadRequest(c, err)
			return
		}

		dummy, err := s.DummyService.CreateDummy(c, &dto)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusCreated, *dummy)
	}
}

// handleDummyGet godoc
// @Summary Gets a dummy.
// @Tags Dummy
// @Description Gets a dummy.
// @Param id path string true "Dummy ID"
// @Accept json
// @Produce json
// @Success 201 {object} services.Dummy
// @Failure 404
// @Failure 500
// @Router /dummies/{id} [get]
func (s *Server) handleDummyGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		dummy, err := s.DummyService.GetDummy(c, id)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "Dummy not found")
			return
		}

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, *dummy)
	}
}

// handleDummyUpdate godoc
// @Summary Updates an existing dummy.
// @Tags Dummy
// @Description Updates an existing dummy.
// @Param id path string true "Dummy ID"
// @Param data body services.DummyDto true "Request body"
// @Accept json
// @Produce json
// @Success 200 {object} services.Dummy
// @Failure 400
// @Failure 500
// @Router /dummies/{id} [post]
func (s *Server) handleDummyUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto services.DummyDto

		if err := c.ShouldBindJSON(&dto); err != nil {
			handleBadRequest(c, err)
			return
		}

		id := c.Param("id")
		updatedDummy, err := s.DummyService.UpdateDummy(c, &dto, id)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusOK, *updatedDummy)
	}
}

// handleDummyDelete godoc
// @Summary Deletes an existing dummy.
// @Tags Dummy
// @Description Deletes an existing dummy.
// @Param id path string true "Dummy ID"
// @Accept json
// @Success 204
// @Failure 500
// @Router /dummies/{id} [delete]
func (s *Server) handleDummyDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := s.DummyService.DeleteDummy(c, id)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusNoContent, "Dummy successfully deleted")
	}
}
