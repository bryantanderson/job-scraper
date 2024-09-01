package services_test

import "github.com/bryantanderson/go-job-assessor/internal/services"

type FakeEventService struct {

}

func InitializeFakeEventService() *FakeEventService {
	return &FakeEventService{}
}

func (s *FakeEventService) Publish(event *services.Event) {

}

func (s *FakeEventService) Subscribe(topic, subscriber string, mChan chan []byte) {

}

func (s *FakeEventService) Register(callback func()) {
	
}


