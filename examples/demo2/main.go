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

	ticketResponse, err := z.Ticketing().Tickets().Show(ctx, 2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(ticketResponse.Ticket.ID)
}
