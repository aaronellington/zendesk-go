package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aaronellington/zendesk-go"
)

func main() {
	ctx := context.Background()

	z := zendesk.NewService(
		os.Getenv("ZENDESK_DEMO_SUBDOMAIN"),
		zendesk.AuthenticationToken{
			Email: os.Getenv("ZENDESK_DEMO_EMAIL"),
			Token: os.Getenv("ZENDESK_DEMO_TOKEN"),
		},
		zendesk.ChatCredentials{
			ClientID:     os.Getenv("ZENDESK_DEMO_CHAT_CLIENT_ID"),
			ClientSecret: os.Getenv("ZENDESK_DEMO_CHAT_CLIENT_SECRET"),
		},
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
	)

	if err := z.Chat().Chats().IncrementalExport(ctx, (time.Now()).Add(time.Hour*-5).Unix(), func(response zendesk.ChatsIncrementalExportResponse) error {
		log.Printf("Found %d", len(response.Chats))
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
