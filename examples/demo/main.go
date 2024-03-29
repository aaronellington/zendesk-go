package main

import (
	"context"
	"log"
	"os"

	"github.com/aaronellington/zendesk-go/zendesk"
)

func main() {
	ctx := context.Background()

	z := zendesk.New(
		os.Getenv("ZENDESK_DEMO_SUBDOMAIN"),
		zendesk.AuthenticationToken{
			Email: os.Getenv("ZENDESK_DEMO_EMAIL"),
			Token: os.Getenv("ZENDESK_DEMO_TOKEN"),
		},
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
	)

	ticketResponse, err := z.Ticketing().AccountSettings().Update(ctx, zendesk.AccountSettingsPayload{
		Settings: map[string]any{
			"branding": map[string]any{
				"header_color": "50C878",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(ticketResponse.Settings.Branding.HeaderColor)
}
