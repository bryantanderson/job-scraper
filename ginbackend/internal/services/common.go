package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
