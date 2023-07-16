package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

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

	defaultStartTime := time.Now().Add(time.Hour * -72)

	check := func(ctx context.Context) {
		if err := z.Chat().AgentsService().UpdateAgentStates(ctx, defaultStartTime); err != nil {
			log.Println(err)
			return
		}

		agentList := z.Chat().AgentsService().GetAgentStates(ctx)

		jsonBytes, _ := json.MarshalIndent(agentList, "", "\t")
		log.Printf("Agent On Chat: %d %s", len(agentList), string(jsonBytes))
		agentList[123] = zendesk.AgentState{}
	}

	check(ctx)
	for range time.NewTicker(time.Second * 5).C {
		check(ctx)
	}
}
