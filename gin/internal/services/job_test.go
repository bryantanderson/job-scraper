package services_test

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/bryantanderson/go-job-assessor/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type JobStoreErrorBehaviorConfig struct {
	errorOnCreate bool
	errorOnQuery  bool
}

type FakeJobStore struct {
	mu   sync.Mutex
	jobs map[string]*services.Job

	errorBehaviorConfig *JobStoreErrorBehaviorConfig
}

func initializeFakeJobStore(cfg *JobStoreErrorBehaviorConfig) *FakeJobStore {
	if cfg == nil {
		cfg = &JobStoreErrorBehaviorConfig{}
	}
	return &FakeJobStore{
		jobs:                make(map[string]*services.Job),
		errorBehaviorConfig: cfg,
	}
}

func (s *FakeJobStore) Create(j *services.Job) error {
	if s.errorBehaviorConfig.errorOnCreate {
		return errors.New(INTERNAL_SERVER_ERROR)
	}
	if _, ok := s.jobs[j.Id]; ok {
		return errors.New("Job already exists")
	}
	s.mu.Lock()
	s.jobs[j.Id] = j
	s.mu.Unlock()
	return nil
}

func (s *FakeJobStore) Query(params map[string]string) ([]*services.Job, error) {
	if s.errorBehaviorConfig.errorOnQuery {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}
	result := make([]*services.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		r := reflect.ValueOf(j)
		for k, v := range params {
			jobValue := reflect.Indirect(r).FieldByName(k)
			if jobValue.String() == v {
				result = append(result, j)
			}
		}
	}
	return result, nil
}

func TestCompleteScrapedJob(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := initializeFakeJobStore(nil)
	sut := services.InitializeJobService(
		InitializeFakeLlmService(),
		fakeStore,
	)

	scrapedJob := services.ScrapedJob{
		Title:       "Senior Software Engineer",
		Company:     "Google",
		Location:    "Mountain View, CA",
		Description: "Works on distributed systems and leads a team of software engineers.",
	}
	completedJob, err := sut.CompleteScrapedJob(&scrapedJob)

	assert.Nil(t, err)
	assert.NotNil(t, completedJob)
	assert.Equal(t, scrapedJob.Title, completedJob.Title)
	assert.Equal(t, scrapedJob.Company, completedJob.Company)
	assert.Equal(t, scrapedJob.Location, completedJob.Location)
	assert.Equal(t, scrapedJob.Description, completedJob.Description)
}

func TestCompleteScrapedJobWithStoreError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	errConfig := &JobStoreErrorBehaviorConfig{
		errorOnCreate: true,
	}
	fakeStore := initializeFakeJobStore(errConfig)
	sut := services.InitializeJobService(
		InitializeFakeLlmService(),
		fakeStore,
	)

	completedJob, err := sut.CompleteScrapedJob(&services.ScrapedJob{})

	assert.Error(t, err)
	assert.Nil(t, completedJob)
}

func TestCompleteScrapedJobWithLlmError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := initializeFakeJobStore(nil)
	sut := services.InitializeJobService(
		InitializeFakeLlmServiceWithError(),
		fakeStore,
	)

	completedJob, err := sut.CompleteScrapedJob(&services.ScrapedJob{})

	assert.Error(t, err)
	assert.Nil(t, completedJob)
}

func TestQueryJobs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := initializeFakeJobStore(nil)
	sut := services.InitializeJobService(
		InitializeFakeLlmService(),
		fakeStore,
	)

	params := make(map[string]string)
	jobs, err := sut.QueryJobs(params)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs))
}
