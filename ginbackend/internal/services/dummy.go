package services

import (
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

func InitializeDummyService(store DummyStore) *DummyService {
	return &DummyService{
		store: store,
	}
}

func (s *DummyService) CreateDummy(dto *DummyDto) (*Dummy, error) {
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

func (s *DummyService) GetDummy(id string) (*Dummy, error) {
	dummy, err := s.store.FindById(id)
	if err != nil {
		return nil, err
	}
	return dummy, nil
}

func (s *DummyService) UpdateDummy(dto *DummyDto, id string) (*Dummy, error) {
	updatedDummy, err := s.store.Update(dto, id)
	if err != nil {
		return nil, err
	}
	return updatedDummy, nil
}

func (s *DummyService) DeleteDummy(id string) error {
	err := s.store.DeleteById(id)
	return err
}
