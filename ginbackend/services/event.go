package services

import (
	"context"

	"sincidium/linkd/api/setup"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	log "github.com/sirupsen/logrus"
)

type Event struct {
	topic       string
	body        string
	contentType string
}

type EventService struct {
	client        *azservicebus.Client
	producerMap   map[string]*azservicebus.Sender
	subscriberMap map[string]*azservicebus.Receiver
	registered    []func()
}

func NewEventService(appSettings *setup.ApplicationSettings) *EventService {
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

	return &EventService{
		client:        client,
		producerMap:   producerMap,
		subscriberMap: subscriberMap,
	}
}

func (s *EventService) Register(callback func()) {
	s.registered = append(s.registered, callback)
}

func (s *EventService) Start() {
	for _, r := range s.registered {
		go r()
	}
}

func (s *EventService) Subscribe(topic, subscriber string, mChan chan []byte) {
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

func (s *EventService) Publish(event *Event) {
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

func (s *EventService) getSubscriber(topic, subscriber string) *azservicebus.Receiver {
	receiver, ok := s.subscriberMap[topic]
	if ok {
		return receiver
	}
	return s.createSubscriber(topic, subscriber)
}

func (s *EventService) createSubscriber(topic, subscriber string) *azservicebus.Receiver {
	receiver, err := s.client.NewReceiverForSubscription(
		topic, subscriber, nil,
	)
	s.subscriberMap[topic] = receiver

	checkFatalError(err)

	return receiver
}

func (s *EventService) Close() {
	for _, p := range s.producerMap {
		p.Close(context.TODO())
	}
	for _, sb := range s.subscriberMap {
		sb.Close(context.TODO())
	}
	s.client.Close(context.TODO())
}
