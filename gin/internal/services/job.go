package services

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	Company           string           `json:"company" jsonschema:"description=The company hiring for the job."`
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
	store  JobStore
	client LlmService
}

func InitializeJobService(
	c LlmService,
	s JobStore,
) *JobService {
	return &JobService{
		store:  s,
		client: c,
	}
}

func (s *JobService) CompleteScrapedJob(scrapedJob *ScrapedJob) (*Job, error) {
	scrapedJobJson, err := json.Marshal(scrapedJob)

	if err != nil {
		log.Errorf("Failed to convert scraped job to JSON: %s", err.Error())
		return nil, err
	}

	prompt := fmt.Sprintf(`
	Given the incomplete job description,
	complete the missing fields "qualifications", "responsibilities", "locationType" and "YearsOfExperience".
    Job Description:
	%s
	`, scrapedJobJson)

	job := &Job{}

	err = s.client.Message(prompt, 500, job)

	if err != nil {
		log.Errorf("Failed to complete scraped job: %s", err.Error())
		return nil, err
	}

	s.copyScrapedJobFields(scrapedJob, job)
	job.Id = uuid.NewString()
	err = s.store.Create(job)

	if err != nil {
		log.Errorf("Failed to store job in database: %s", err.Error())
		return nil, err
	}

	return job, nil
}

func (s *JobService) QueryJobs(params map[string]string) ([]*Job, error) {
	return s.store.Query(params)
}

func (s *JobService) copyScrapedJobFields(scrapedJob *ScrapedJob, jobToComplete *Job) {
	jobToComplete.Title = scrapedJob.Title
	jobToComplete.Company = scrapedJob.Company
	jobToComplete.Description = scrapedJob.Description
	jobToComplete.Location = scrapedJob.Location
}
