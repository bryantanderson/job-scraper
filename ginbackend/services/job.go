package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sincidium/linkd/api/setup"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
)

const (
	TIMEOUT int32 = 3600
)

type Qualification struct {
	Description string `json:"description"`
}

type Responsibility struct {
	Description string `json:"description"`
}

type Job struct {
	Id                string           `json:"id,omitempty" bson:"_id"`
	Title             string           `json:"title" jsonschema:"description=The position that the job is hiring for"`
	Description       string           `json:"description" jsonschema:"description=Brief summary of what the job is about"`
	Responsibilities  []Responsibility `json:"responsibilities" jsonschema:"description=Responsibilities of the ideal employee listed out in the job description"`
	Qualifications    []Qualification  `json:"qualifications" jsonschema:"description=Required qualifications of the ideal employee"`
	Location          string           `json:"location" jsonschema:"description=The location that the job requires candidate to be in, set as NaN if not specified or remote"`
	LocationType      string           `json:"locationType" jsonschema:"description=The location type for the job. Choose one of [onsite, hybrid, remote]."`
	YearsOfExperience int8             `json:"yearsOfExperience" jsonschema:"description=The number of years of experience for the ideal employee"`
	ElasticId         string           `json:"elasticId,omitempty"`
}

type JobStore interface {
	Create(j *Job) error
	Query(params map[string]string) ([]*Job, error)
}

type JobService struct {
	inTopic             string
	inTopicSubscription string
	outTopic            string
	store               JobStore
	client              *instructor.InstructorOpenAI
	eventService        *EventService
}

func InitializeJobService(
	s *setup.ApplicationSettings,
	c *instructor.InstructorOpenAI,
	e *EventService,
	js JobStore,
) *JobService {
	srv := &JobService{
		inTopic:             s.JobTasksTopic,
		inTopicSubscription: topicNameToSubscriptionName(s.JobTasksTopic),
		outTopic:            s.JobResultsTopic,
		store:               js,
		client:              c,
		eventService:        e,
	}
	srv.RegisterSubscribers()
	return srv
}

func (s *JobService) RegisterSubscribers() {
	routine := func() {
		numWorkers := 1
		mChan := make(chan []byte, numWorkers)
		defer close(mChan)
		for i := 1; i <= numWorkers; i++ {
			go s.worker(i, mChan)
		}
		s.eventService.Subscribe(s.inTopic, s.inTopicSubscription, mChan)
	}
	s.eventService.Register(routine)
}

func (s *JobService) worker(i int, mChan <-chan []byte) {
	log.Printf("Worker number %d for job tasks starting...", i)
	for {
		body, ok := <-mChan
		description := string(body)

		if !ok {
			// If the channel is closed, the worker should stop
			break
		}

		s.structureJob(description)
	}
}

func (s *JobService) structureJob(description string) {
	prompt := fmt.Sprintf(`
	Given the job description, turn it into a structured json format.

    Job Description:
	%s
	`, description)

	job := Job{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(TIMEOUT))
	defer cancel()

	resp, err := s.client.CreateChatCompletion(
		ctx,
		makeChatCompletionRequest(prompt, len(description)/2),
		&job,
	)
	_ = resp

	if err != nil {
		log.Errorf("Failed to structure job description: %s", err.Error())
		return
	}

	jobJson, err := json.Marshal(job)

	if err != nil {
		log.Errorln(err)
		return
	}

	event := Event{
		topic:       s.outTopic,
		body:        string(jobJson),
		contentType: "application/json",
	}
	s.eventService.Publish(&event)
}

func (s *JobService) CompleteScrapedJob(scrapedJob *ScrapedJob) *Job {
	scrapedJobJson, err := json.Marshal(scrapedJob)

	if err != nil {
		log.Errorf("Failed to convert scraped job to JSON: %s", err.Error())
		return nil
	}

	prompt := fmt.Sprintf(`
	Given the incomplete job description,
	complete the missing fields "qualifications", "responsibilities", "locationType" and "YearsOfExperience".
    Job Description:
	%s
	`, scrapedJobJson)

	job := &Job{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(TIMEOUT))
	defer cancel()

	resp, err := s.client.CreateChatCompletion(
		ctx,
		makeChatCompletionRequest(prompt, 0),
		job,
	)
	_ = resp

	if err != nil {
		log.Errorf("Failed to complete scraped job: %s", err.Error())
		return nil
	}

	job.Id = uuid.NewString()
	err = s.store.Create(job)

	if err != nil {
		log.Errorf("Failed to store job in database: %s", err.Error())
		return nil
	}

	return job
}

func (s *JobService) QueryJobs(params map[string]string) ([]*Job, error) {
	return s.store.Query(params)
}
