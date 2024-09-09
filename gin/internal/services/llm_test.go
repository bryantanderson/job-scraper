package services_test

import (
	"errors"

	"github.com/bryantanderson/go-job-assessor/internal/services"
)

type ErrorHistory struct {
	jobCriteria            bool
	assessRequirements     bool
	assessCompatibility    bool
	assessLocation         bool
	assessResponsibilities bool
	assessSkills           bool
}

type FakeLlmService struct {
}

type FakeLlmServiceWithError struct {
}

func InitializeFakeLlmService() *FakeLlmService {
	return &FakeLlmService{}
}

func InitializeFakeLlmServiceWithError() *FakeLlmServiceWithError {
	return &FakeLlmServiceWithError{}
}

func (s *FakeLlmService) Message(prompt string, maxTokens int, resType any) error {
	switch t := resType.(type) {
	case *services.Score:
		t.Explanation = "Great Score"
		t.Score = 100
	case *services.Match:
		t.IsMatch = true
	case *services.Point:
		t.Explanation = "Great Point"
		t.IsValid = true
	}
	return nil
}

func (s *FakeLlmServiceWithError) Message(prompt string, maxTokens int, resType any) error {
	return errors.New(INTERNAL_SERVER_ERROR)
}
