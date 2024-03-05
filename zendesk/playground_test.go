package zendesk_test

// import (
// 	"context"
// 	"errors"
// 	"net"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/aaronellington/zendesk-go/zendesk"
// 	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
// 	"github.com/gobwas/ws"
// 	"github.com/gobwas/ws/wsutil"
// )

// func NewWebSocketTestingThing() (*WebSocketTestingThing, chan error) {
// 	errorChan := make(chan error)
// 	return &WebSocketTestingThing{
// 		error:  errorChan,
// 		Marker: make(chan int),
// 	}, make(chan error)
// }

// type WebSocketTestingThing struct {
// 	MaxPongResponses int
// 	ActiveToken      string
// 	Marker           chan int
// 	error            chan error
// }

// func (w *WebSocketTestingThing) Build() *zendesk.Service {
// 	client, server := net.Pipe()

// 	serverWriter := wsutil.NewWriter(server, ws.StateServerSide, ws.OpText)
// 	serverReader := wsutil.NewReader(server, ws.StateServerSide)

// 	// Handle incoming messages from the client (our service)
// 	go func() {
// 		// Set up a limit on the expected messages
// 		messageCount := 0
// 		pongsSent := 0
// 		for {
// 			header, err := serverReader.NextFrame()
// 			if err != nil {
// 				w.error <- err
// 			}

// 			if header.OpCode == ws.OpPing && pongsSent != w.MaxPongResponses {
// 				pongsSent += 1
// 				serverWriter.Reset(server, ws.StateServerSide, ws.OpPong)

// 				if err := ws.WriteHeader(serverWriter, header); err != nil {
// 					w.error <- err
// 				}

// 				_, err := serverWriter.Write(nil)
// 				if err != nil {
// 					w.error <- err
// 				}

// 				if err := serverWriter.Flush(); err != nil {
// 					w.error <- err
// 				}

// 				messageCount++
// 			}

// 			if err := serverReader.Discard(); err != nil {
// 				w.error <- err
// 			}

// 			// Once the limit on messages is reached, trigger the "shutdown" of the server
// 			if messageCount >= 2 {
// 				w.Marker <- messageCount
// 				w.error <- nil
// 			}
// 		}

// 	}()

// 	return zendesk.NewService(
// 		"example",
// 		zendesk.AuthenticationToken{
// 			Email: "example@example.com",
// 			Token: "abc123",
// 		},
// 		zendesk.ChatCredentials{
// 			ClientID:     "fake-client-id",
// 			ClientSecret: "fake-client-secret",
// 		},
// 		zendesk.WithRoundTripper(study.RoundTripper{
// 			RoundTripFunc: func(r *http.Request) (*http.Response, error) {
// 				return &http.Response{}, nil
// 			},
// 		}),
// 		zendesk.WithWebsocketConnection(&client),
// 	)
// }

// func TestPlayground(t *testing.T) {
// 	ctx := context.Background()

// 	testThing, errChan := NewWebSocketTestingThing()
// 	testThing.MaxPongResponses = -1

// 	z := testThing.Build()

// 	go func() {
// 		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
// 			if errors.Is(err, context.Canceled) {
// 				return
// 			}

// 			errChan <- err
// 		}
// 	}()

// 	select {
// 	case <-testThing.Marker:
// 		break
// 	case err := <-errChan:
// 		t.Fatal(err)
// 	}

// 	actualTimeSinceLastFrame := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetTimeSinceLastFrameSent()
// 	if actualTimeSinceLastFrame == nil {
// 		t.Fatalf("expected to have recorded ping")
// 	}

// 	if *actualTimeSinceLastFrame > time.Second*10 {
// 		t.Fatal("expected a ping within 10 seconds")
// 	}
// }
