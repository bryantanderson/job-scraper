package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/bryantanderson/go-job-assessor/docs"
	"github.com/bryantanderson/go-job-assessor/internal/database"
	"github.com/bryantanderson/go-job-assessor/internal/handlers"
	"github.com/bryantanderson/go-job-assessor/internal/services"
	"github.com/bryantanderson/go-job-assessor/internal/setup"
	log "github.com/sirupsen/logrus"
)

func run(ctx context.Context) {
	// Get app settings
	settings := setup.ReadApplicationSettings()

	// Instantiate server
	server := handlers.NewServer(settings)
	defer server.Close()

	// Connect to MongoDB
	db := database.InitializeDatabase(settings)
	db.Open()

	// Connect to OpenAI
	openai := services.InitializeLlmService(settings.AzureOpenAiApiKey, settings.AzureOpenAiEndpoint)

	// Connect to Elastic Search
	elasticSearch := database.InitializeElasticSearch(settings)

	// Instantiate event service
	eventService := services.InitializeEventServiceImpl(settings)

	// Instantiate candidate service
	candidateStore := database.InitializeCandidateStore(db)
	candidateService := services.InitializeCandidateService(candidateStore)

	// Instantiate assessor service
	assessStore := database.InitializeAssessStore(db, elasticSearch)
	assessorService := services.InitializeAssessorService(
		settings.AssessmentTasksTopic,
		settings.AssessmentResultsTopic,
		openai,
		eventService,
		assessStore,
	)

	// Instantiate job service
	jobStore := database.InitializeJobStore(db, elasticSearch)
	jobService := services.InitializeJobService(
		openai,
		jobStore,
	)

	// Instantiate scraper service
	scraperService := services.InitializeScraperService(jobService, assessorService)

	server.AddRoutes(
		db,
		elasticSearch,
		jobService,
		eventService,
		scraperService,
		assessorService,
		candidateService,
	)

	httpServer := &http.Server{
		Addr:         server.Settings.ServerPort,
		Handler:      server.Router,
		ReadTimeout:  server.Settings.ServerReadTimeout * time.Second,
		WriteTimeout: server.Settings.ServerWriteTimeout * time.Second,
	}
	server.EventService.Start()

	log.Infof("Listening on %s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err == http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Error listening and serving: %s\n", err)
	}
}

func main() {
	ctx := context.Background()
	run(ctx)
}
