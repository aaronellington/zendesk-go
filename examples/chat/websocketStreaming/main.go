package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
)

func main() {

	ctx, _ := context.WithCancel(context.Background())

	// initalGoroutineCount := runtime.NumGoroutine()

	defer func() {

		time.Sleep(time.Millisecond * 250)
		// finalGoroutineCount := runtime.NumGoroutine()
		// if finalGoroutineCount > initalGoroutineCount {
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
		// }

	}()

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

	go z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAgentMetric(ctx, zendesk.LiveChatMetricKeyAgentsOnline)

	if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
		log.Printf("Websocket exiting. Here is the error message: %s", err.Error())
	}
}
