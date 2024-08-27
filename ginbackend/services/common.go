package services

import (
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

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

func checkFatalError(err error) {
	if err != nil {
		panic(err)
	}
}

func convertToJson(toConvert any) (string, error) {
	jsonBytes, err := json.Marshal(toConvert)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), err
}

func topicNameToSubscriptionName(t string) string {
	return fmt.Sprintf("sbs-%s", t)
}
