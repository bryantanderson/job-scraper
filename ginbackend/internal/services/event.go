package services

import (
	"context"
	"sync"

	"time"

	"github.com/bryantanderson/go-job-assessor/internal/setup"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	log "github.com/sirupsen/logrus"
)

type Event struct {
	topic       string
	body        string
	contentType string
}

type EventService interface {
	Publish(event *Event)
	Subscribe(topic, subscriber string, mChan chan []byte)
	Register(callback func())
}

type EventServiceImpl struct {
	mu            sync.Mutex
	client        *azservicebus.Client
	producerMap   map[string]*azservicebus.Sender
	subscriberMap map[string]*azservicebus.Receiver
	registered    []func()
}

func InitializeEventServiceImpl(appSettings *setup.ApplicationSettings) *EventServiceImpl {
	producerMap := make(map[string]*azservicebus.Sender)
	subscriberMap := make(map[string]*azservicebus.Receiver)
	client, err := azservicebus.NewClientFromConnectionString(appSettings.ServiceBusConnectionString, nil)

	checkFatalError(err)

	// Register producer for dispatching finished job related tasks
	jobSender, err := client.NewSender(appSettings.JobTasksTopic, nil)
	producerMap[appSettings.JobTasksTopic] = jobSender

	checkFatalError(err)

	// Register producer for dispatching finished assessment related tasks
	assessSender, err := client.NewSender(appSettings.AssessmentTasksTopic, nil)
	producerMap[appSettings.AssessmentTasksTopic] = assessSender

	checkFatalError(err)

	return &EventServiceImpl{
		client:        client,
		producerMap:   producerMap,
		subscriberMap: subscriberMap,
	}
}

func (s *EventServiceImpl) Register(callback func()) {
	s.registered = append(s.registered, callback)
}

func (s *EventServiceImpl) Start() {
	for _, r := range s.registered {
		go r()
	}
}

func (s *EventServiceImpl) Publish(event *Event) {
	producer, ok := s.producerMap[event.topic]
	if !ok {
		log.Errorf("Producer does not exist for topic %s\n", event.topic)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	message := azservicebus.Message{
		Body:        []byte(event.body),
		ContentType: &event.contentType,
	}
	err := producer.SendMessage(ctx, &message, nil)

	if err != nil {
		log.Errorf("Unable to produce message: %s\n", err.Error())
	}
}

func (s *EventServiceImpl) Subscribe(topic, subscriber string, mChan chan []byte) {
	receiver := s.getSubscriber(topic, subscriber)
	for {
		messages, err := receiver.ReceiveMessages(context.TODO(), 10, nil)
		if err != nil {
			log.Errorf("Failed to receive Messages: %s\n", err.Error())
		}
		for _, m := range messages {
			body := m.Body

			if body == nil {
				log.Errorf("Failed to parse message body: %s\n", body)
			}

			log.Infoln("Successfully received message")

			mChan <- body
			err = receiver.CompleteMessage(context.TODO(), m, nil)

			if err != nil {
				log.Errorf("Failed to complete message: %s\n", err.Error())
			}
		}
	}
}

func (s *EventServiceImpl) getSubscriber(topic, subscriber string) *azservicebus.Receiver {
	receiver, ok := s.subscriberMap[topic]
	if ok {
		return receiver
	}
	return s.createSubscriber(topic, subscriber)
}

func (s *EventServiceImpl) createSubscriber(topic, subscriber string) *azservicebus.Receiver {
	receiver, err := s.client.NewReceiverForSubscription(
		topic, subscriber, nil,
	)
	s.mu.Lock()
	s.subscriberMap[topic] = receiver
	s.mu.Unlock()

	checkFatalError(err)

	return receiver
}

func (s *EventServiceImpl) Close() {
	for _, p := range s.producerMap {
		p.Close(context.TODO())
	}
	for _, sb := range s.subscriberMap {
		sb.Close(context.TODO())
	}
	s.client.Close(context.TODO())
}
