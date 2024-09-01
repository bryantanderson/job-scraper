package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"github.com/bryantanderson/go-job-assessor/internal/database"
	_ "github.com/bryantanderson/go-job-assessor/docs"
	"github.com/bryantanderson/go-job-assessor/internal/handlers"
	"github.com/bryantanderson/go-job-assessor/internal/setup"
	"time"
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
	openai := setup.InitializeOpenAI(settings)

	// Connect to Elastic Search
	elasticSearch := database.InitializeElasticSearch(settings)

	server.AddRoutes(db, openai, elasticSearch)

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
