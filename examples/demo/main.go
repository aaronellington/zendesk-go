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

	newTagsAfterAdd, err := z.Support().Tickets().AddTags(ctx, 6170, zendesk.Tags{
		"foobar",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", newTagsAfterAdd)

	newTagsAfterRemoval, err := z.Support().Tickets().RemoveTags(ctx, 6170, zendesk.Tags{
		"foobar",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", newTagsAfterRemoval)

	newTagsAfterSet, err := z.Support().Tickets().SetTags(ctx, 6170, zendesk.Tags{
		"foobar",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", newTagsAfterSet)
}
