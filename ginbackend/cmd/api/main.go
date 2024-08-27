package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sincidium/linkd/api/database"
	"sincidium/linkd/api/handlers"
	"sincidium/linkd/api/setup"
	_ "sincidium/linkd/api/docs"
	"time"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

func run(ctx context.Context) {
	// Get app settings
	settings := setup.ReadApplicationSettings()

	// Instantiate server
	server := handlers.NewServer(settings)
	defer server.Close()

	// Connect to MongoDB
	db := database.InitDB(settings)
	db.Open()

	// Connect to OpenAI
	openai := openai.NewClient(settings.OpenAiApiKey)
	client := instructor.FromOpenAI(
		openai,
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)

	// Connect to Elastic Search
	elasticSearch := database.InitElasticSearch(settings)

	server.AddRoutes(db, client, elasticSearch)

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
