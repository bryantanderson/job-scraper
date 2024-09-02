package services_test

import (
	"errors"
)

type FakeLlmService struct {
}

type FakeLlmServiceWithError struct {
}

func InitializeFakeLlmService() *FakeLlmService {
	return &FakeLlmService{}
}

func (s *FakeLlmService) Message(prompt string, maxTokens int, resType any) error {
	return nil
}

func (s *FakeLlmServiceWithError) Message(prompt string, maxTokens int, resType any) error {
	return errors.New("Internal Server Error")
}
