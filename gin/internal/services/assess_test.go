package services_test

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/bryantanderson/go-job-assessor/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AssessorStoreErrorBehaviorConfig struct {
	errorOnCreate                    bool
	errorOnCreateInternalJobCriteria bool
	errorOnQueryInternalJobCriteria  bool
	errorOnFindById                  bool
	errorOnQuery                     bool
	errorOnDelete                    bool
}

type FakeAssessorStore struct {
	mu           sync.Mutex
	assessments  map[string]*services.Assessment
	jobCriterion map[string]*services.JobCriteria

	// For controlling error behavior
	errorBehaviorConfig *AssessorStoreErrorBehaviorConfig
}

func InitializeFakeAssessorStore(cfg *AssessorStoreErrorBehaviorConfig) *FakeAssessorStore {
	if cfg == nil {
		cfg = &AssessorStoreErrorBehaviorConfig{}
	}
	return &FakeAssessorStore{
		assessments:         make(map[string]*services.Assessment),
		jobCriterion:        make(map[string]*services.JobCriteria),
		errorBehaviorConfig: cfg,
	}
}

func (s *FakeAssessorStore) Create(a *services.Assessment) error {
	if s.errorBehaviorConfig.errorOnCreate {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	if _, ok := s.assessments[a.Id]; ok {
		return errors.New("Assessment already exists")
	}
	s.mu.Lock()
	s.assessments[a.Id] = a
	s.mu.Unlock()
	return nil
}

func (s *FakeAssessorStore) CreateInternalJobCriteria(jc *services.JobCriteria) error {
	if s.errorBehaviorConfig.errorOnCreateInternalJobCriteria {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	if _, ok := s.jobCriterion[jc.Id]; ok {
		return errors.New("Job Criteria already exists")
	}
	s.mu.Lock()
	s.jobCriterion[jc.Id] = jc
	s.mu.Unlock()
	return nil
}

func (s *FakeAssessorStore) QueryInternalJobCriteria(id string) (*services.JobCriteria, error) {
	if s.errorBehaviorConfig.errorOnQueryInternalJobCriteria {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}
	if jc, ok := s.jobCriterion[id]; ok {
		return jc, nil
	}
	return nil, errors.New("Job Criteria does not exist")
}

func (s *FakeAssessorStore) FindById(assessmentId string) (*services.Assessment, error) {
	if s.errorBehaviorConfig.errorOnFindById {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}

	if a, ok := s.assessments[assessmentId]; ok {
		return a, nil
	}
	return nil, errors.New("Assessment does not exist")
}

func (s *FakeAssessorStore) Query(params map[string]string) ([]*services.Assessment, error) {
	if s.errorBehaviorConfig.errorOnQuery {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}
	result := make([]*services.Assessment, 0, len(s.assessments))
	for _, a := range s.assessments {
		r := reflect.ValueOf(a)
		// Check all params
		for k, v := range params {
			assessmentValue := reflect.Indirect(r).FieldByName(k)
			if assessmentValue.String() == v {
				result = append(result, a)
			}
		}
	}
	return result, nil
}

func (s *FakeAssessorStore) Delete(userId string) error {
	if s.errorBehaviorConfig.errorOnDelete {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	if _, ok := s.assessments[services.UserIdToAssessmentId(userId)]; ok {
		s.mu.Lock()
		delete(s.assessments, userId)
		s.mu.Unlock()
		return nil
	}
	return errors.New("Assessment does not exist")
}

/*
 *
 */

func initializeTestAssessorService(store services.AssessorStore, llm services.LlmService, event services.EventService) *services.AssessorService {
	if store == nil {
		store = InitializeFakeAssessorStore(nil)
	}
	if llm == nil {
		llm = InitializeFakeLlmService()
	}
	if event == nil {
		event = InitializeFakeEventService()
	}
	return services.InitializeAssessorService("test", "test", llm, event, store)
}

func TestAssessCandidateAndGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currentDate := time.Now()

	testCases := []struct {
		name string
		data services.AssessPayload
	}{
		{
			name: "Empty job location",
			data: services.AssessPayload{
				Job:    services.Job{},
				UserId: "testUserOne",
				Candidate: services.Candidate{
					Experiences: []services.Experience{
						services.Experience{
							Title:       "Director of Engineering",
							Company:     "Google",
							Description: "Built Google from scratch",
							StartDate:   currentDate,
							EndDate:     nil,
						},
						services.Experience{
							Title:       "Senior Software Engineer",
							Company:     "NASA",
							Description: "Built NASA from scratch",
							StartDate:   time.Now().Add(time.Duration(-2400) * time.Hour),
							EndDate:     &currentDate,
						},
					},
				},
			},
		},
		{
			name: "Matching job and candidate location",
			data: services.AssessPayload{
				Job: services.Job{
					Location: "San Francisco",
				},
				UserId: "testUserTwo",
				Candidate: services.Candidate{
					Location: "San Francisco",
				},
			},
		},
		{
			name: "Different job and candidate locations",
			data: services.AssessPayload{
				Job: services.Job{
					Location: "San Francisco",
				},
				UserId: "testUserThree",
				Candidate: services.Candidate{
					Location: "Palo Alto",
				},
			},
		},
	}
	sut := initializeTestAssessorService(nil, nil, nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assessment, err := sut.AssessCandidate(&tc.data)
			assert.NotNil(t, assessment)
			assert.NoError(t, err)
			// Test that the assessment was indeed stored
			existingAssessment, err := sut.GetAssessment(tc.data.UserId)
			assert.NoError(t, err)
			assert.NotNil(t, existingAssessment)
		})
	}
}

func TestAssessCandidateWithLlmError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := services.AssessPayload{
		UserId: "testUserOne",
	}
	fakeStore := InitializeFakeAssessorStore(nil)
	fakeLlmService := InitializeFakeLlmServiceWithError()
	sut := initializeTestAssessorService(fakeStore, fakeLlmService, nil)

	sut.AssessCandidate(&payload)
	existingAssessment, err := sut.GetAssessment(payload.UserId)

	assert.Error(t, err)
	assert.Nil(t, existingAssessment)
}

func TestAssessCandidateWithStoreError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name   string
		config *AssessorStoreErrorBehaviorConfig
	}{
		{
			name:   "Errors while creating a job criteria",
			config: &AssessorStoreErrorBehaviorConfig{errorOnCreateInternalJobCriteria: true},
		},
		{
			name:   "Errors while creating an assessment",
			config: &AssessorStoreErrorBehaviorConfig{errorOnCreate: true},
		},
	}
	payload := services.AssessPayload{
		UserId: "testUserOne",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeStore := InitializeFakeAssessorStore(tc.config)
			fakeLlmService := InitializeFakeLlmService()
			sut := initializeTestAssessorService(fakeStore, fakeLlmService, nil)

			sut.AssessCandidate(&payload)
			existingAssessment, err := sut.GetAssessment(payload.UserId)

			assert.Error(t, err)
			assert.Nil(t, existingAssessment)
		})
	}
}

func TestQueryAssessments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sut := initializeTestAssessorService(nil, nil, nil)

	query := make(map[string]string)
	query["id"] = "test"

	existingAssessments, err := sut.QueryAssessments(query)

	assert.NoError(t, err)
	assert.Equal(t, len(existingAssessments), 0)
}
