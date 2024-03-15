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
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
	)

	// NOTE: This is fine to do before initiating a connection. The library will wait up to 15 seconds for a connection to be established, and then perform any queued writes
	go z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAgentMetricByDepartment(ctx, zendesk.LiveChatMetricKeyAgentsOnline, 13388700431505)

	// NOTE: Connecting to the WebSocket will consume frames from the Zendesk API until an error occurs. It also handles checking for a stale connection and sending keepalive messages
	// to the Zendesk Server.
	if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
		log.Printf("Websocket exiting. Here is the error message: %s", err.Error())
	}
}
