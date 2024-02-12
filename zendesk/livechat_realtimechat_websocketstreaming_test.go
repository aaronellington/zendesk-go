package zendesk_test

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	client, server := net.Pipe()

	z := createTestWebsocketService(
		t,
		[]study.RoundTripFunc{
			createSuccessfulChatAuth(t),
			study.ServeAndValidate(
				t,
				&study.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/responses/livechat/realtimechat/get_all_chat_metrics_200.json",
				},
				study.ExpectedTestRequest{
					Method: http.MethodGet,
					Path:   "/stream/chats",
				},
			),
		},
		&client,
	)

	serverWriter := wsutil.NewWriter(server, ws.StateServerSide, ws.OpText)
	serverReader := wsutil.NewReader(server, ws.StateServerSide)

	testError := make(chan error)
	testMarker := make(chan int)

	// Handle incoming messages from the client (our service)
	go func() error {
		// Set up a limit on the expected messages
		messageCount := 0
		for {
			header, err := serverReader.NextFrame()
			if err != nil {
				testError <- err
			}

			if header.OpCode == ws.OpPing {

				serverWriter.Reset(server, ws.StateServerSide, ws.OpPong)

				if err := ws.WriteHeader(serverWriter, header); err != nil {
					testError <- err
				}

				_, err := serverWriter.Write(nil)
				if err != nil {
					testError <- err
				}

				if err := serverWriter.Flush(); err != nil {
					testError <- err
				}

				messageCount++
			}

			if err := serverReader.Discard(); err != nil {
				testError <- err
			}

			// Once the limit on messages is reached, trigger the "shutdown" of the server
			if messageCount >= 2 {
				testMarker <- messageCount
				testError <- nil
			}
		}

	}()

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	// Wait for the messageCount to be reached before progressing
	<-testMarker

	actualTimeSinceLastFrame := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetTimeSinceLastFrameReceived()
	if actualTimeSinceLastFrame == nil {
		t.Fatalf("expected to have recorded ping")
	}

	if *actualTimeSinceLastFrame > time.Second*10 {
		t.Fatal("expected a ping within 10 seconds")
	}

	if err := <-testError; err != nil {
		t.Fatal(err)
	}

}
