package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
)

func prettyPrint(v any) error {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	os.Stdout.Write(bytes)
	os.Stdout.WriteString("\n")

	return nil
}

func PrintErr(err error) {
	if err == nil {
		return
	}

	zdErr, ok := err.(*zendesk.Error)
	if ok {
		log.Fatalf("Zendesk Error: [%d] %s", zdErr.StatusCode, string(zdErr.Body))
	}

	log.Fatal(err)
}

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

	err := z.LiveChat().AgentEvent().IncrementalExport(ctx, time.Now().Add(time.Minute*-8), func(response zendesk.AgentEventExportResponse) error {
		_ = prettyPrint(response)
		return nil
	})
	if err != nil {
		PrintErr(err)
	}
}
