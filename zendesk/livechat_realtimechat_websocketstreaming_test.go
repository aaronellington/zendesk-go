package zendesk_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate that the Upgrade Header  is being sent, as well as Auth header

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			testError <- err
		}
		go func() {
			defer conn.Close()

			var (
				state        = ws.StateServerSide
				serverReader = wsutil.NewReader(conn, state)
				serverWriter = wsutil.NewWriter(conn, state, ws.OpText)
			)

			for {
				header, err := serverReader.NextFrame()
				if err != nil {
					testError <- err
				}

				if header.OpCode == ws.OpPing {

					serverWriter.Reset(conn, ws.StateServerSide, ws.OpPong)

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

				}

				if err := serverReader.Discard(); err != nil {
					testError <- err
				}

				testError <- nil
			}
		}()
	}))

	rtcWSHost := strings.TrimPrefix(mockServer.URL, "http")

	z := createTestRealTimeChatWebsocketService(
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
		fmt.Sprintf("ws%s", rtcWSHost),
	)

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("Websocket exiting, restarting. Here is the error message: %s", err.Error())
			}
		}
	}()

	time.Sleep(time.Second * 6)

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
