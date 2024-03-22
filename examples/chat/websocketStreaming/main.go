package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	// NOTE: This is fine to do before initiating a connection. The library will wait up to 15 seconds for a connection to the Websocket to be established, and then perform any queued writes
	go func() {
		departments, err := z.LiveChat().Chat().Department().List(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, department := range departments {
			if !department.Enabled {
				continue
			}

			if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllAgentMetrics(&department.ID); err != nil {
				log.Fatal(err)
			}

			if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllChatMetrics(&department.ID); err != nil {
				log.Fatal(err)
			}
		}

		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllAgentMetrics(nil); err != nil {
			log.Fatal(err)
		}
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAllChatMetrics(nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			log.Printf("Websocket exiting. Here is the error message: %s", err.Error())
		}
	}()

	// Every 15 seconds pull the state from the Internal Memory Cache and prettyPrint this to the console.
	ticker := time.NewTicker(time.Second * 15)
	for t := range ticker.C {
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

		fmt.Println(t.String())
		fmt.Println("=== Chat  Data ===")
		prettyPrint(chat)
		fmt.Println("=== Agent Data ===")
		prettyPrint(agent)
		fmt.Println("=== ========== ===")
	}
}
