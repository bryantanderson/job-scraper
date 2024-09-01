package services_test

import (
	"errors"
	"sincidium/linkd/api/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type FakeCandidateStore struct {
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
	if candidate, ok := s.candidates[id]; ok {
		return candidate, nil
	}
	return nil, errors.New("Candidate not found")
}

func (s *FakeCandidateStore) Delete(id string) error {
	if _, ok := s.candidates[id]; ok {
		delete(s.candidates, id)
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

func TestDeleteCandidate(t *testing.T) {
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

	err = sut.DeleteCandidate(createdCandidate.Id)

	assert.NoError(t, err)

	err = sut.DeleteCandidate(createdCandidate.Id)

	assert.Error(t, err)
}
