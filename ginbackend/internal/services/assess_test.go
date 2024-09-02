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

type FakeAssessorStore struct {
	mu           sync.Mutex
	assessments  map[string]*services.Assessment
	jobCriterion map[string]*services.JobCriteria
}

func InitializeFakeAssessorStore() *FakeAssessorStore {
	return &FakeAssessorStore{
		assessments:  make(map[string]*services.Assessment),
		jobCriterion: make(map[string]*services.JobCriteria),
	}
}

func (s *FakeAssessorStore) Create(a *services.Assessment) error {
	if _, ok := s.assessments[a.Id]; ok {
		return errors.New("Assessment already exists")
	}
	s.mu.Lock()
	s.assessments[a.Id] = a
	s.mu.Unlock()
	return nil
}

func (s *FakeAssessorStore) CreateInternalJobCriteria(jc *services.JobCriteria) error {
	if _, ok := s.jobCriterion[jc.Id]; ok {
		return errors.New("Job Criteria already exists")
	}
	s.mu.Lock()
	s.jobCriterion[jc.Id] = jc
	s.mu.Unlock()
	return nil
}

func (s *FakeAssessorStore) QueryInternalJobCriteria(id string) (*services.JobCriteria, error) {
	if jc, ok := s.jobCriterion[id]; ok {
		return jc, nil
	}
	return nil, errors.New("Job Criteria does not exist")
}

func (s *FakeAssessorStore) FindById(assessmentId string) (*services.Assessment, error) {
	if a, ok := s.assessments[assessmentId]; ok {
		return a, nil
	}
	return nil, errors.New("Assessment does not exist")
}

func (s *FakeAssessorStore) Query(params map[string]string) ([]*services.Assessment, error) {
	result := make([]*services.Assessment, 0, len(s.assessments))
	for _, a := range s.assessments {
		r := reflect.ValueOf(a)
		// Check all params
		for k, v := range params {
			assessmentValue := reflect.Indirect(r).FieldByName(k)
			if assessmentValue.String() == v {
				result = append(result, a)
				return result, nil
			}
		}
	}
	return result, nil
}

func (s *FakeAssessorStore) Delete(userId string) error {
	if _, ok := s.assessments[services.UserIdToAssessmentId(userId)]; ok {
		s.mu.Lock()
		delete(s.assessments, userId)
		s.mu.Unlock()
		return nil
	}
	return errors.New("Assessment does not exist")
}

func initializeTestAssessorService(t *testing.T) *services.AssessorService {
	t.Helper()
	fakeStore := InitializeFakeAssessorStore()
	fakeLlmService := InitializeFakeLlmService()
	fakeEventService := InitializeFakeEventService()
	return services.InitializeAssessorService("test", "test", fakeLlmService, fakeEventService, fakeStore)
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
				Job:       services.Job{},
				UserId:    "testUserOne",
				Candidate: services.Candidate{
					Experiences: []services.Experience{
						services.Experience{
							Title: "Director of Engineering",
							Company: "Google",
							Description: "Built Google from scratch",
							StartDate: currentDate,
							EndDate: nil,
							IsCurrent: true,
						},
						services.Experience{
							Title: "Senior Software Engineer",
							Company: "NASA",
							Description: "Built NASA from scratch",
							StartDate: time.Now().Add(time.Duration(-2400) * time.Hour),
							EndDate: &currentDate,
							IsCurrent: false,
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
	sut := initializeTestAssessorService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut.AssessCandidate(&tc.data)
			existingAssessment, err := sut.GetAssessment(tc.data.UserId)
			assert.NoError(t, err)
			assert.NotNil(t, existingAssessment)
		})
	}
}

func TestQueryAssessments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sut := initializeTestAssessorService(t)

	query := make(map[string]string)
	query["id"] = "test"

	existingAssessments, err := sut.QueryAssessments(query)

	assert.NoError(t, err)
	assert.Equal(t, len(existingAssessments), 0)
}
