# zendesk-go
Zendesk API client library for Go.

## Getting Started
### Install
```shell
go get github.com/aaronellington/zendesk-go
```

## Create new connection
```go
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
        // Logger is optional, see implementing to see how to add your custom logger here
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
        // Optionally set http.RoundTripper - this is helpful when writing tests
		zendesk.WithRoundTripper(customRoundTripper),
	)

	tags, err := support.Tickets().AddTags(ctx, 6170, zendesk.Tags{
		"foobar",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", tags)
}
```
