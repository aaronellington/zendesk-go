package zendesk_test

import (
	"context"
	"errors"
	"fmt"
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

	// This is our test mockserver
	mockableZendeskRTCWebsocketServer := &mockRealTimeChatWebsocketServer{
		State: State{},
		Settings: Settings{
			ValidOAuthToken: "Bearer fake-token",
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != mockableZendeskRTCWebsocketServer.Settings.ValidOAuthToken {
			testError <- fmt.Errorf("expected '%s', got: '%s'", mockableZendeskRTCWebsocketServer.Settings.ValidOAuthToken, authorizationHeader)

			w.WriteHeader(http.StatusForbidden)
			return
		}

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			testError <- err
			return
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
					return
				}

				if header.OpCode == ws.OpPing {

					serverWriter.Reset(conn, ws.StateServerSide, ws.OpPong)

					if err := ws.WriteHeader(serverWriter, header); err != nil {
						testError <- err
						return
					}

					_, err := serverWriter.Write(nil)
					if err != nil {
						testError <- err
						return
					}

					if err := serverWriter.Flush(); err != nil {
						testError <- err
						return
					}

				}

				if err := serverReader.Discard(); err != nil {
					testError <- err
					return
				}

				testError <- nil
			}
		}()
	}))
	defer mockServer.Close()

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
				testError <- err
				return
			}
		}
	}()

	time.Sleep(time.Second * 6)

	go func() {
		actualTimeSinceLastFrame := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetTimeSinceLastFrameReceived()
		if actualTimeSinceLastFrame == nil {
			testError <- fmt.Errorf("expected to have recorded ping")
			return
		}

		if *actualTimeSinceLastFrame > time.Second*10 {
			testError <- fmt.Errorf("expected to received ping within 10 seconds")
		}

	}()

	if err := <-testError; err != nil {
		t.Fatal(err)
	}
}

type mockRealTimeChatWebsocketServer struct {
	State    State
	Settings Settings
}

type State struct {
	lastWriteFromClient *time.Time
}

type Settings struct {
	ValidOAuthToken string
	WriteTimeout    time.Duration
}
