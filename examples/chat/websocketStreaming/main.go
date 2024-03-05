package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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
	c1 := runtime.NumGoroutine()

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

	_ = ctx
	_ = z

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			log.Printf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
		}
	}()

	time.Sleep(time.Second * 20)

	c2 := runtime.NumGoroutine()
	if c2 > c1 {
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)

		log.Printf("Error too many goroutines: %d extra", c2-c1)
	}
}
