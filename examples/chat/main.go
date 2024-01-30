package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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
		errBytes, _ := io.ReadAll(zdErr.Response.Body)
		log.Fatalf("Zendesk Error: [%d] %s", zdErr.Response.StatusCode, string(errBytes))
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

	departments, err := z.LiveChat().Chat().Department().List(ctx)
	for _, department := range departments {
		if !department.Settings.ChatEnabled {
			continue
		}
		deptID := department.ID
		fmt.Println(len(departments))
		fmt.Println("Processing department ID:", deptID)
		metricResponse, err := z.LiveChat().RealTimeChat().REST().GetChatMetricsByDepartment(ctx, deptID)
		_ = prettyPrint(metricResponse)

		if err != nil {
			PrintErr(err)
		}
	}

	if err != nil {
		PrintErr(err)
	}

}
