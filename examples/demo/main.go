package main

import (
	"context"
	"log"
	"os"

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
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
	)

	if err := z.Guide().Articles().List(ctx, func(response zendesk.ArticlesResponse) error {
		log.Printf("Found %d articles", len(response.Articles))
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
