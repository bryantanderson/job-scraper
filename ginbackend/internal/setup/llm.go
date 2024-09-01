package setup

import (
	"context"
	"time"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	instructor *instructor.InstructorOpenAI
}

func InitializeOpenAI(settings *ApplicationSettings) *OpenAI {
	config := openai.DefaultAzureConfig(settings.AzureOpenAiApiKey, settings.AzureOpenAiEndpoint)
	openai := openai.NewClientWithConfig(config)
	instructorClient := instructor.FromOpenAI(
		openai,
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)
	return &OpenAI{
		instructor: instructorClient,
	}
}

func (o *OpenAI) Message(prompt string, maxTokens int, resType any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	resp, err := o.instructor.CreateChatCompletion(
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