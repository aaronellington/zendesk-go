# zendesk-go

<!-- aaronellington/stencil -->

[![Go](https://github.com/aaronellington/zendesk-go/actions/workflows/go.yml/badge.svg)](https://github.com/aaronellington/zendesk-go/actions/workflows/go.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/aaronellington/zendesk-go.svg)](https://pkg.go.dev/github.com/aaronellington/zendesk-go) [![Go Report Card](https://goreportcard.com/badge/github.com/aaronellington/zendesk-go)](https://goreportcard.com/report/github.com/aaronellington/zendesk-go)

<!-- aaronellington/stencil -->

![zendesk-go logo](./ops/images/zendesk-go.png)

Zendesk API client library for Go.

## Getting Started

### Install

```shell
go get github.com/aaronellington/zendesk-go
```

### Create new connection

> [!NOTE]
> You will need to set your Zendesk Credentials in your environment before the below example will work. In bash, you can do this by running:  
> `export ZENDESK_DEMO_EMAIL=<YOUR_EMAIL>; export ZENDESK_DEMO_TOKEN=<YOUR_TOKEN>; bash;`

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/aaronellington/zendesk-go/zendesk"
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
		// Logger is optional, see implementation to see how to add your custom logger here
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
