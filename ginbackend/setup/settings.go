package setup

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ApplicationSettings struct {
	DatabaseName string
	DatabaseUri  string

	ServiceBusNamespace        string
	JobTasksTopic              string
	JobResultsTopic            string
	AssessmentTasksTopic       string
	AssessmentResultsTopic     string
	ServiceBusConnectionString string

	AzureOpenAiEndpoint string
	AzureOpenAiApiKey   string

	ElasticCloudId string
	ElasticApiKey  string

	ServerPort         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
}

func ReadApplicationSettings() *ApplicationSettings {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	settings := ApplicationSettings{}

	settings.ServerPort = os.Getenv("GIN_ADDR")

	ginReadTimeout, err := strconv.Atoi(os.Getenv("GIN_READ_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	ginWriteTimeout, err := strconv.Atoi(os.Getenv("GIN_WRITE_TIMEOUT"))
	if err != nil {
		panic(err)
	}

	settings.ServiceBusConnectionString = os.Getenv("SERVICE_BUS_CONNECTION_STRING")

	// Gin service is a receiver for job related processing
	settings.JobTasksTopic = os.Getenv("JOB_TASKS_TOPIC")

	// Gin service publishes job related processing results
	settings.JobResultsTopic = os.Getenv("JOB_RESULTS_TOPIC")

	// Gin service is a receiver for assessment related processing
	settings.AssessmentTasksTopic = os.Getenv("ASSESSMENT_TASKS_TOPIC")

	// Gin service publishes assessment related processing results
	settings.AssessmentResultsTopic = os.Getenv("ASSESSMENT_RESULTS_TOPIC")

	settings.ServerReadTimeout = time.Duration(ginReadTimeout)
	settings.ServerWriteTimeout = time.Duration(ginWriteTimeout)
	settings.DatabaseName = os.Getenv("DATABASE_NAME")
	settings.DatabaseUri = os.Getenv("DATABASE_CONNECTION_STRING")

	settings.AzureOpenAiEndpoint = os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	settings.AzureOpenAiApiKey = os.Getenv("AZURE_OPEN_AI_API_KEY")

	settings.ElasticCloudId = os.Getenv("ELASTIC_CLOUD_ID")
	settings.ElasticApiKey = os.Getenv("ELASTIC_API_KEY")

	return &settings
}
