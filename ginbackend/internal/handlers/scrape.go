package handlers

import (
	"errors"
	"net/http"

	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerScrapeRoutes() {
	scrapeGroup := s.Router.Group("/scrape")
	{
		scrapeGroup.POST("/seek", s.handleSeekScrape())
		scrapeGroup.POST("/indeed", s.handleIndeedScrape())
		scrapeGroup.GET("/seek/:user", s.handleSeekScrapeGet())
	}
}

// handleSeekScrape godoc
// @Summary Scrapes a seek job page
// @Tags Scrape
// @Description Scrapes a seek job page. Accepts a "candidate" and ranks the candidate against the job postings if the flag is selected.
// @Param data body services.ScrapePayload true "Request body"
// @Accept json
// @Produce json
// @Success 200 {array} services.ScrapedJob
// @Failure 400
// @Router /scrape/seek [post]
func (s *Server) handleSeekScrape() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload services.ScrapePayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			handleBadRequest(c, err)
			return
		}

		scrapedJobs := s.ScraperService.ScrapeSeekJobPage(&payload)
		c.JSON(http.StatusOK, scrapedJobs)
	}
}

// handleIndeedScrape godoc
// @Summary Scrapes an indeed job page
// @Tags Scrape
// @Description Scrapes an Indeed job page. Accepts a "candidate" and ranks the candidate against the job postings if the flag is selected.
// @Param data body services.ScrapeIndeedPayload true "Request body"
// @Accept json
// @Produce json
// @Success 200 {array} services.ScrapedJob
// @Failure 400
// @Router /scrape/indeed [post]
func (s *Server) handleIndeedScrape() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload services.ScrapeIndeedPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			handleBadRequest(c, err)
			return
		}

		scrapedJobs := s.ScraperService.ScrapeIndeedJobPage(&payload)
		c.JSON(http.StatusOK, scrapedJobs)
	}
}

// handleSeekScrapeGet godoc
// @Summary Fetches the assessments made for a particular candidate, in relation to scraped jobs.
// @Tags Scrape
// @Description Fetches the assessments made for a particular candidate, in relation to scraped jobs.
// @Param user path string true "User ID"
// @Accept json
// @Produce json
// @Success 200 {array} services.ScrapedJobAssessment
// @Failure 400
// @Failure 500
// @Router /scrape/seek/{user} [get]
func (s *Server) handleSeekScrapeGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user")

		if userId == "" {
			handleBadRequest(c, errors.New("Invalid user"))
			return
		}

		assessments := s.ScraperService.GetAssessments(userId)

		if assessments == nil {
			handleInternalError(c, errors.New("Internal server error"))
			return
		}

		c.JSON(http.StatusOK, assessments)
	}
}
