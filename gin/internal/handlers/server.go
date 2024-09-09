package handlers

import (
	"github.com/bryantanderson/go-job-assessor/internal/database"
	"github.com/bryantanderson/go-job-assessor/internal/services"
	"github.com/bryantanderson/go-job-assessor/internal/setup"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router   *gin.Engine
	Settings *setup.ApplicationSettings

	EventService     *services.EventServiceImpl
	JobService       *services.JobService
	AssessorService  *services.AssessorService
	ScraperService   *services.ScraperService
	CandidateService *services.CandidateService
}

func NewServer(settings *setup.ApplicationSettings) *Server {
	return &Server{
		Router:   gin.Default(),
		Settings: settings,
	}
}

func (s *Server) AddRoutes(
	db *database.Database,
	elastic *database.ElasticDatabase,
	jobService *services.JobService,
	eventService *services.EventServiceImpl,
	scraperService *services.ScraperService,
	assessorService *services.AssessorService,
	candidateService *services.CandidateService,
) {
	// Define middleware
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())

	s.registerAssessRoutes()
	s.registerHealthRoutes()
	s.registerScrapeRoutes()
	s.registerInstrumentRoutes()

	// Define Swagger route
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	s.EventService = eventService
	s.JobService = jobService
	s.CandidateService = candidateService
	s.AssessorService = assessorService
	s.ScraperService = scraperService
}

func (s *Server) Close() {
	s.EventService.Close()
}
