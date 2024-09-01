package services_test

import (
	"errors"
	"testing"

	"github.com/bryantanderson/go-job-assessor/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Define the fake

type FakeDummyStore struct {
	dummies map[string]*services.Dummy
}

func InitializeFakeDummyStore() *FakeDummyStore {
	return &FakeDummyStore{
		dummies: make(map[string]*services.Dummy),
	}
}

func (s *FakeDummyStore) Create(dummy *services.Dummy) error {
	if _, ok := s.dummies[dummy.Id]; ok {
		return errors.New("Dummy already exists")
	}
	s.dummies[dummy.Id] = dummy
	return nil
}

func (s *FakeDummyStore) FindById(id string) (*services.Dummy, error) {
	if dummy, ok := s.dummies[id]; ok {
		return dummy, nil
	}
	return nil, errors.New("Dummy not found")
}

func (s *FakeDummyStore) Update(dto *services.DummyDto, id string) (*services.Dummy, error) {
	if dummy, ok := s.dummies[id]; ok {
		dummy.Name = dto.Name
		return dummy, nil
	}
	return nil, errors.New("Dummy not found")
}

func (s *FakeDummyStore) DeleteById(id string) error {
	if _, ok := s.dummies[id]; ok {
		delete(s.dummies, id)
		return nil
	}
	return errors.New("Dummy not found")
}

func TestCreateAndGetDummy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeDummyStore()
	sut := services.InitializeDummyService(fakeStore)

	dto := &services.DummyDto{Name: "Test"}

	// Create dummy and make assertions
	createdDummy, err := sut.CreateDummy(dto)

	assert.NoError(t, err)
	assert.NotNil(t, createdDummy)
	assert.Equal(t, dto.Name, createdDummy.Name)
	assert.NotEmpty(t, createdDummy.Id)

	// Find dummy and make assertions
	storedDummy, err := sut.GetDummy(createdDummy.Id)

	assert.NoError(t, err)
	assert.Equal(t, createdDummy.Name, storedDummy.Name)
}

func TestUpdateDummy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeDummyStore()
	sut := services.InitializeDummyService(fakeStore)

	dto := &services.DummyDto{Name: "Test"}

	// Create dummy
	createdDummy, err := sut.CreateDummy(dto)

	assert.NoError(t, err)

	// Update dummy with new data and make assertions
	updateDto := &services.DummyDto{Name: "Updated Test"}
	updatedDummy, err := sut.UpdateDummy(updateDto, createdDummy.Id)

	assert.NoError(t, err)
	assert.Equal(t, updateDto.Name, updatedDummy.Name)
}

func TestDeleteDummy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStore := InitializeFakeDummyStore()
	sut := services.InitializeDummyService(fakeStore)

	dto := &services.DummyDto{Name: "Test"}

	// Create dummy
	createdDummy, err := sut.CreateDummy(dto)

	assert.NoError(t, err)

	// Delete dummy and make assertions
	err = sut.DeleteDummy(createdDummy.Id)

	assert.NoError(t, err)

	// Check that the dummy no longer exists
	existingDummy, err := sut.GetDummy(createdDummy.Id)

	assert.Error(t, err)
	assert.Nil(t, existingDummy)
}
