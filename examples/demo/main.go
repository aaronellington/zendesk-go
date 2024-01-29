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

	err := z.LiveChat().ChatDepartments().List(ctx, func(departments []zendesk.Department) error {
		for _, department := range departments {
			deptID := department.ID
			fmt.Println("Processing department ID:", deptID)

			err := z.LiveChat().ChatStream().List(ctx, fmt.Sprint(deptID), func(response zendesk.ChatsStreamResponse) error {
				fmt.Println("Completed processing for department ID:", deptID)
				_ = prettyPrint(response)
				return nil
			})

			if err != nil {
				PrintErr(err)
			}
		}
		return nil
	})
	if err != nil {
		PrintErr(err)
	}

}
