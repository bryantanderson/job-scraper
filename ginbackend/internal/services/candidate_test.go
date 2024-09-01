package services_test

import (
	"errors"
	"testing"

	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type FakeCandidateStore struct {
	candidates map[string]*services.Candidate
}

type FakeCandidateStoreWithError struct {
	candidates map[string]*services.Candidate
}

func InitializeFakeCandidateStore() *FakeCandidateStore {
	return &FakeCandidateStore{
		candidates: make(map[string]*services.Candidate),
	}
}

func (s *FakeCandidateStore) Create(c *services.Candidate) error {
	if _, ok := s.candidates[c.Id]; ok {
		return errors.New("Candidate already exists")
	}
	s.candidates[c.Id] = c
	return nil
}

func (s *FakeCandidateStore) FindById(id string) (*services.Candidate, error) {
	return findCandidateById(id, s.candidates)
}

func (s *FakeCandidateStore) Delete(id string) error {
	return deleteCandidate(id, s.candidates)
}

/*
 *
 * The below receivers are for the fake store used to test error scenarios within the service.
 *
 */

func InitializeFakeCandidateStoreWithError() *FakeCandidateStoreWithError {
	return &FakeCandidateStoreWithError{
		candidates: make(map[string]*services.Candidate),
	}
}

func (s *FakeCandidateStoreWithError) Create(c *services.Candidate) error {
	return errors.New("Internal Server Error")
}

func (s *FakeCandidateStoreWithError) FindById(id string) (*services.Candidate, error) {
	return findCandidateById(id, s.candidates)
}

func (s *FakeCandidateStoreWithError) Delete(id string) error {
	return deleteCandidate(id, s.candidates)
}

/*
 *
 * Base methods to be called for both fake stores
 *
 */

func findCandidateById(id string, candidates map[string]*services.Candidate) (*services.Candidate, error) {
	if candidate, ok := candidates[id]; ok {
		return candidate, nil
	}
	return nil, errors.New("Candidate not found")
}

func deleteCandidate(id string, candidates map[string]*services.Candidate) error {
	if _, ok := candidates[id]; ok {
		delete(candidates, id)
		return nil
	}
	return errors.New("Dummy not found")
}

func TestCreateAndGetCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeCandidateStore()
	sut := services.InitializeCandidateService(fakeStore)

	dto := services.CandidateDto{
		Summary:  "This is a unit test",
		Location: "San Francisco, California",
	}
	createdCandidate, err := sut.CreateCandidate(&dto)

	assert.NoError(t, err)
	assert.NotNil(t, createdCandidate)
	assert.Equal(t, dto.Summary, createdCandidate.Summary)
	assert.Equal(t, dto.Location, createdCandidate.Location)
	assert.NotEmpty(t, createdCandidate.Id)

	existingCandidate, err := sut.GetCandidate(createdCandidate.Id)

	assert.NoError(t, err)
	assert.NotNil(t, existingCandidate)
	assert.Equal(t, createdCandidate.Id, existingCandidate.Id)
}

func TestCreateCandidateWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeCandidateStoreWithError()
	sut := services.InitializeCandidateService(fakeStore)

	dto := services.CandidateDto{
		Summary:  "This is a unit test",
		Location: "San Francisco, California",
	}
	candidate, err := sut.CreateCandidate(&dto)

	assert.Error(t, err)
	assert.Nil(t, candidate)
}

func TestGetNonExistentCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeCandidateStore()
	sut := services.InitializeCandidateService(fakeStore)

	candidate, err := sut.GetCandidate("Not a real ID")

	assert.Error(t, err)
	assert.Nil(t, candidate)
}

func TestDeleteCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeCandidateStore()
	sut := services.InitializeCandidateService(fakeStore)

	dto := services.CandidateDto{
		Skills: []string{
			"Python",
			"Golang",
			"Java",
			"TypeScript",
		},
		Summary:  "This is a unit test",
		Location: "San Francisco, California",
	}

	createdCandidate, err := sut.CreateCandidate(&dto)

	assert.NoError(t, err)
	assert.NotNil(t, createdCandidate)

	err = sut.DeleteCandidate(createdCandidate.Id)

	assert.NoError(t, err)

	err = sut.DeleteCandidate(createdCandidate.Id)

	assert.Error(t, err)
}
