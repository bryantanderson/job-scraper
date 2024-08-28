package handlers

import (
	"errors"
	"net/http"
	"sincidium/linkd/api/services"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerScrapeRoutes() {
	scrapeGroup := s.Router.Group("/scrape")
	{
		scrapeGroup.POST("/seek", s.handleSeekScrape())
		scrapeGroup.GET("/seek/:user", s.handleSeekScrapeGet())
	}
}

// handleSeekScrape godoc
// @Summary Scrapes a seek job page 
// @Tags Scrape
// @Description Scrapes a seek job page for Software Engineering Jobs. Accepts a "candidate" and ranks the candidate against the job postings.
// @Param data body services.ScrapeSeekPayload true "Request payload"
// @Accept json
// @Produce json
// @Success 200
// @Router /scrape/seek [post]
func (s *Server) handleSeekScrape() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload services.ScrapeSeekPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			handleBadRequest(c, err)
			return
		}

		scrapedJobs := s.ScraperService.ScrapeSeek(&payload)
		c.JSON(http.StatusOK, makeData(scrapedJobs))
	}
}

// handleSeekScrapeGet godoc
// @Summary Fetches the assessments made for a particular candidate, in relation to scraped jobs.
// @Tags Scrape
// @Description Fetches the assessments made for a particular candidate, in relation to scraped jobs.
// @Param user path string true "User ID"
// @Accept json
// @Produce json
// @Success 200
// @Router /scrape/seek/{user} [get]
func (s *Server) handleSeekScrapeGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user")
		
		if userId == "" {
			handleBadRequest(c, errors.New("Invalid user"))
			return
		}

		assessments := s.ScraperService.GetScrapedSeekAssessments(userId)
		
		if assessments == nil {
			handleInternalError(c, errors.New("Internal server error"))
			return
		}

		c.JSON(http.StatusOK, makeData(assessments))
	}
}