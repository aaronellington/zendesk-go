package main

import (
	"context"
	"log"
	"os"

	"github.com/aaronellington/zendesk-go"
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

	if err := z.Support().Organizations().IncrementalExport(ctx, 0, func(response zendesk.OrganizationsIncrementalExportResponse) error {
		log.Printf("Found %d Organizations: %d", len(response.Organizations), response.EndTime)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := z.Support().Users().IncrementalExport(ctx, 0, func(response zendesk.UsersIncrementalExportResponse) error {
		log.Printf("Found %d Users: %d", len(response.Users), response.EndTime)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := z.Support().Tickets().IncrementalExport(ctx, 0, 500, func(response zendesk.TicketsIncrementalExportResponse) error {
		log.Printf("Found %d Tickets: %d", len(response.Tickets), response.EndTime)

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
