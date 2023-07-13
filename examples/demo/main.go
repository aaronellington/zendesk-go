package main

import (
	"context"
	"fmt"
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

	query := fmt.Sprintf("timestamp:[%s TO *]", time.Now().UTC().Add(time.Hour*-3).Format(time.RFC3339))

	if err := z.Chat().Chats().Search(ctx, query, func(page zendesk.ChatsSearchResponse) error {
		for _, searchResult := range page.Results {
			chat, err := z.Chat().Chats().Show(ctx, searchResult.ID)
			if err != nil {
				return err
			}

			log.Printf("Bla: %s", chat.EndTimestamp)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// if err := z.Chat().Chats().IncrementalExport(ctx, (time.Now()).Add(time.Hour*-5).Unix(), func(response zendesk.ChatsIncrementalExportResponse) error {
	// 	for _, chat := range response.Chats {
	// 		for _, eng := range chat.ChatEngagements {
	// 			log.Printf("eng: %s", eng.ID)
	// 		}
	// 	}

	// 	return nil
	// }); err != nil {
	// 	log.Fatal(err)
	// }
}
