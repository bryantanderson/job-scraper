package services

import (
	"encoding/json"
	"fmt"
	"strings"

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

func handleGoroutineError(err error, errChan chan<- error) {
	// If the error is caused by a cancel, we do not need to propagate the error back up
	if strings.Contains(err.Error(), "context canceled") {
		return
	}
	errChan <- err
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
