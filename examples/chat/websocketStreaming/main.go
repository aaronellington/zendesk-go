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

	// departmentID := zendesk.GroupID(13388700431505)
	timeWindow := zendesk.LiveChatTimeWindow30Minutes

	// NOTE: This is fine to do before initiating a connection. The library will wait up to 15 seconds for a connection to be established, and then perform any queued writes
	departments, err := z.LiveChat().Chat().Department().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, department := range departments {
			if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllAgentMetrics(&department.ID); err != nil {
				log.Fatal(err)
			}
		}
		go z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllAgentMetrics(nil)
		go z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToSingleChatWindowMetric(zendesk.LiveChatMetricKeyMissedChats, &timeWindow, nil)
		go z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToSingleChatMetric(zendesk.LiveChatMetricKeyIncomingChats, nil)
	}()

	// NOTE: Connecting to the WebSocket will consume frames from the Zendesk API until an error occurs. It also handles checking for a stale connection and sending keepalive messages
	// to the Zendesk Server.
	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			log.Printf("Websocket exiting. Here is the error message: %s", err.Error())
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		log.Println("==Report==")
		agent, err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetAllAgentMetricsForDepartments()
		if err != nil {
			log.Println(err)
			continue
		}

		chat, err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetAllChatMetricsForDepartments()
		if err != nil {
			log.Println(err)
			continue
		}

		prettyPrint(chat[0].MissedChats.ThirtyMinuteWindow)
		prettyPrint(agent)
	}
}
