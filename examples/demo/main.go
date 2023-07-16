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

	type AgentStatus struct {
		Status          zendesk.AgentEventValue
		EngagementCount zendesk.AgentEventValue
	}

	agentList := map[zendesk.UserID]AgentStatus{}
	startTime := time.Now().Add(time.Hour * -100)

	check := func(ctx context.Context) {
		log.Println("Checking....")
		if err := z.Chat().AgentsService().IncrementalExport(ctx, startTime, func(response zendesk.AgentEventExportResponse) error {
			for _, agentTimeline := range response.AgentEvents {
				existingAgent, alreadyExists := agentList[agentTimeline.AgentID]
				if !alreadyExists {
					existingAgent = AgentStatus{
						Status:          "",
						EngagementCount: "0",
					}
				}

				switch agentTimeline.FieldName {
				case "engagements":
					existingAgent.EngagementCount = agentTimeline.Value
				case "status":
					existingAgent.Status = agentTimeline.Value

				}
				if agentTimeline.Value == "offline" {
					delete(agentList, agentTimeline.AgentID)
					continue
				}

				agentList[agentTimeline.AgentID] = existingAgent
			}

			startTime = response.EndTime()

			return nil
		}); err != nil {
			log.Println(err)
		}

		jsonBytes, _ := json.MarshalIndent(agentList, "", "\t")
		log.Printf("Agent On Chat: %d %s", len(agentList), string(jsonBytes))
	}

	check(ctx)
	for range time.NewTicker(time.Second * 5).C {
		check(ctx)
	}
}
