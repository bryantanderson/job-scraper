package services

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Dummy struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DummyDto struct {
	Name string `json:"name"`
}

type DummyStore interface {
	Create(dummy *Dummy) error
	FindById(id string) (*Dummy, error)
	Update(dto *DummyDto, id string) (*Dummy, error)
	DeleteById(id string) error
}

type DummyService struct {
	store DummyStore
}

func NewDummyService(store DummyStore) *DummyService {
	return &DummyService{
		store: store,
	}
}

func (s *DummyService) CreateDummy(c *gin.Context, dto *DummyDto) (*Dummy, error) {
	dummy := Dummy{
		Id:   uuid.NewString(),
		Name: dto.Name,
	}
	err := s.store.Create(&dummy)
	if err != nil {
		return nil, err
	}
	return &dummy, nil
}

func (s *DummyService) GetDummy(c *gin.Context, id string) (*Dummy, error) {
	dummy, err := s.store.FindById(id)
	if err != nil {
		return nil, err
	}
	return dummy, nil
}

func (s *DummyService) UpdateDummy(c *gin.Context, dto *DummyDto, id string) (*Dummy, error) {
	updatedDummy, err := s.store.Update(dto, id)
	if err != nil {
		return nil, err
	}
	return updatedDummy, nil
}

func (s *DummyService) DeleteDummy(c *gin.Context, id string) error {
	err := s.store.DeleteById(id)
	return err
}
