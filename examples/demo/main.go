package main

import (
	"context"
	"encoding/json"
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
		"kaseya1523304719",
		zendesk.AuthenticationToken{
			Email: "zendesk_integrations@kaseya.com",
			Token: "xQ6FrYWjSVyeQHCM3u1vpjhooyqhncXeRrbhM8lb",
		},
		zendesk.ChatCredentials{
			ClientID:     "RLc1Ddddmtq88haQ64rm9ewGj8qgWuFSSSdaKTewllBwURq8N7",
			ClientSecret: "EGfWIklds8kVXV1riZGpanb7W3w9hpkMsN5GOHvMauZtHOGb5KJfxEVtRVJcEl7b",
		},
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk Chat API - ", log.LstdFlags)),
	)

	// err := z.LiveChat().ChatStream().List(ctx, "13388700431505", func(response zendesk.ChatsStreamResponse) error {
	// 	_ = prettyPrint(response)
	// 	return nil
	// })
	state := z.LiveChat().AgentEvent().GetAgentStates(ctx)
	prettyPrint(state)
}
