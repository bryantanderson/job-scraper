package services_test

import "github.com/bryantanderson/go-job-assessor/internal/services"

type FakeAssessorStore struct {
	assessments map[string]*services.Assessment
	rubrics map[string]*services.Rubric
}

func InitializeFakeAssessorStore() *FakeAssessorStore {
	return &FakeAssessorStore{
		assessments: make(map[string]*services.Assessment),
		rubrics: make(map[string]*services.Rubric),
	}
}

func (s *FakeAssessorStore) Create(a *services.Assessment) error {
	return nil
}

func (s *FakeAssessorStore) CreateInternalJobCriteria(jc *services.Rubric) error {
	return nil
}

func (s *FakeAssessorStore) QueryInternalJobCriteria(id string) (*services.Rubric, error) {
	return nil, nil
}

func (s *FakeAssessorStore) FindById(userId string) (*services.Assessment, error) {
	return nil, nil
}

func (s *FakeAssessorStore) Query(params map[string]string) ([]*services.Assessment, error) {
	return nil, nil
}

func (s *FakeAssessorStore) Delete(userId string) error {
	return nil
}

