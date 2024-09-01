package handlers

import (
	"github.com/bryantanderson/go-job-assessor/internal/database"
	"github.com/bryantanderson/go-job-assessor/internal/middleware"
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
	DummyService     *services.DummyService
	JobService       *services.JobService
	AssessorService  *services.AssessorService
	ScraperService   *services.ScraperService
	CandidateService *services.CandidateService
}

func NewServer(settings *setup.ApplicationSettings) *Server {
	s := &Server{
		Router:   gin.Default(),
		Settings: settings,
	}
	return s
}

func (s *Server) AddRoutes(
	db *database.Database,
	client *setup.OpenAI,
	elastic *database.ElasticDatabase,
) {
	// Define middleware
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(middleware.ErrorHandler)

	s.registerAssessRoutes()
	s.registerCandidateRoutes()
	s.registerHealthRoutes()
	s.registerDummyRoutes()
	s.registerScrapeRoutes()
	s.registerInstrumentRoutes()

	// Define Swagger route
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Instantiate event service
	eventService := services.InitializeEventServiceImpl(s.Settings)

	// Instantiate candidate service
	candidateStore := database.InitializeCandidateStore(db)
	candidateService := services.InitializeCandidateService(candidateStore)

	// Instantiate dummy service
	dummyStore := database.InitializeDummyStore(db)
	dummyService := services.InitializeDummyService(dummyStore)

	// Instantiate assessor service
	assessStore := database.InitializeAssessStore(db, elastic)
	assessorService := services.InitializeAssessorService(
		s.Settings.AssessmentTasksTopic, 
		s.Settings.AssessmentResultsTopic, 
		client, 
		eventService, 
		assessStore,
	)

	// Instantiate job service
	jobStore := database.InitializeJobStore(db, elastic)
	jobService := services.InitializeJobService(
		s.Settings.JobTasksTopic,
		s.Settings.JobResultsTopic, 
		client, 
		eventService, 
		jobStore,
	)

	// Instantiate scraper service
	scraperService := services.InitializeScraperService(jobService, assessorService)

	s.EventService = eventService
	s.DummyService = dummyService
	s.JobService = jobService
	s.CandidateService = candidateService
	s.AssessorService = assessorService
	s.ScraperService = scraperService
}

func (s *Server) Close() {
	s.EventService.Close()
}
