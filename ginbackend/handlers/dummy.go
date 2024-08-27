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

		c.JSON(http.StatusCreated, makeData(*dummy))
	}
}

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

		c.JSON(http.StatusOK, makeData(*dummy))
	}
}

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

		c.JSON(http.StatusOK, makeData(*updatedDummy))
	}
}

func (s *Server) handleDummyDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := s.DummyService.DeleteDummy(c, id)

		if err != nil {
			handleInternalError(c, err)
			return
		}

		c.JSON(http.StatusNoContent, makeData("Deleted"))
	}
}
