package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	// subscribed := false

	for range t.C {
		fmt.Println("Time since last frame: ", z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetTimeSinceLastFrameSent())

		// if !subscribed && z.LiveChat().RealTimeChat().RealTimeChatStreamingService().WebsocketReady() {
		// 	if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAgentMetric(ctx, zendesk.LiveChatMetricKeyAgentsOnline); err != nil {
		// 		log.Fatalf(err.Error())
		// 		return
		// 	}

		// 	subscribed = true

		// 	log.Println("subbed to metric!!")
		// }
	}
}
