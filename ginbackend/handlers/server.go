package handlers

import (
	"github.com/gin-gonic/gin"
	"sincidium/linkd/api/services"
	"sincidium/linkd/api/setup"
	"sincidium/linkd/api/database"
	"sincidium/linkd/api/middleware"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router   *gin.Engine
	Settings *setup.ApplicationSettings

	EventService    *services.EventService
	DummyService    *services.DummyService
	JobService      *services.JobService
	AssessorService *services.AssessorService
	ScraperService 	*services.ScraperService
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
	client *instructor.InstructorOpenAI,
	elastic *database.ElasticDatabase,
) {
	// Define middleware
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(middleware.ErrorHandler)

	s.registerAssessRoutes()
	s.registerHealthRoutes()
	s.registerDummyRoutes()
	s.registerScrapeRoutes()
	s.registerInstrumentRoutes()

	// Define Swagger route
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Instantiate event service
	eventService := services.NewEventService(s.Settings)

	// Instantiate dummy service
	dummyStore := database.NewDummyStore(db)
	dummyService := services.NewDummyService(dummyStore)

	// Instantiate assessor service
	assessStore := database.NewAssessStore(db, elastic)
	assessorService := services.NewAssessorService(s.Settings, client, eventService, assessStore)

	// Instantiate job service
	jobStore := database.NewJobStore(db, elastic)
	jobService := services.NewJobService(s.Settings, client, eventService, jobStore)

	// Instantiate scraper service
	scraperService := services.NewScraperService(jobService, assessorService)

	s.EventService = eventService
	s.DummyService = dummyService
	s.JobService = jobService
	s.AssessorService = assessorService
	s.ScraperService = scraperService
}

func (s *Server) Close() {
	s.EventService.Close()
}
