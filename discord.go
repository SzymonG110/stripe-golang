package main

import (
	"github.com/go-zoox/fetch"
	"os"
)

func SendWebhook(text string) (bool, error) {
	discordWebhookURL := os.Getenv("STRIPE_DISCORD_WEBHOOK")

	if !ParseStringToBool(discordWebhookURL) {
		return false, nil
	}

	_, err := fetch.Post(discordWebhookURL, &fetch.Config{
		Body: map[string]interface{}{
			"content": text,
		},
	})

	if err != nil {
		return false, err
	}

	return true, err
}
