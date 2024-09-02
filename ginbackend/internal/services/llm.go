package services

import (
	"context"
	"time"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"
)

type LlmService interface {
	Message(prompt string, maxTokens int, resType any) error
}

type Llm struct {
	instructor *instructor.InstructorOpenAI
}

func InitializeLlmService(apiKey, endpoint string) *Llm {
	config := openai.DefaultAzureConfig(apiKey, endpoint)
	openai := openai.NewClientWithConfig(config)
	instructorClient := instructor.FromOpenAI(
		openai,
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)
	return &Llm{
		instructor: instructorClient,
	}
}

func (l *Llm) Message(prompt string, maxTokens int, resType any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	resp, err := l.instructor.CreateChatCompletion(
		ctx,
		makeChatCompletionRequest(prompt, maxTokens),
		resType,
	)
	_ = resp // sends back original response so no information loss from original API

	return err
}


func makeChatCompletionRequest(prompt string, maxTokens int) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.0,
	}
}