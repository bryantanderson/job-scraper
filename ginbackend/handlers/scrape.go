package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerScrapeRoutes() {
	scrapeGroup := s.Router.Group("/scrape")
	{
		scrapeGroup.POST("/", s.handleScrape())
	}
}

// Scraping godoc
// @Summary Scrapes a given URL
// @Tags Scrape
// @Description Scrapes a given URL
// @Accept json
// @Produce json
// @Success 200
// @Router /scrape [post]
func (s *Server) handleScrape() gin.HandlerFunc {
	return func(c *gin.Context) {
		scrapedJobs := s.ScraperService.ScrapeJobs()
		c.JSON(http.StatusOK, makeData(scrapedJobs))
	}
}