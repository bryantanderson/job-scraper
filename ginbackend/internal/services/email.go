package services

import "github.com/resend/resend-go/v2"

type Email struct {
	client *resend.Client
}

func InitializeEmail(apiKey string) *Email {
	client := resend.NewClient(apiKey)
	return &Email{
		client: client,
	}
}


