package zendesk_test

import (
	"context"
	"encoding/json"
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
		state: State{},
		settings: Settings{
			ValidOAuthToken: "Bearer fakea-token",
		},
		testError: testError,
	}

	mockServer := mockableZendeskRTCWebsocketServer.createDefaultServer()

	mockServer.Start()
	defer mockServer.Close()

	rtcWSHost := strings.TrimPrefix(mockServer.URL, "http")

	z := createTestRealTimeChatWebsocketService(
		t,
		[]study.RoundTripFunc{
			createSuccessfulChatAuth(t),
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

	// go func() {
	// 	actualTimeSinceLastFrame := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetTimeSinceLastFrameReceived()
	// 	if actualTimeSinceLastFrame == nil {
	// 		testError <- fmt.Errorf("expected to have recorded ping")
	// 		return
	// 	}

	// 	if *actualTimeSinceLastFrame > time.Second*10 {
	// 		testError <- fmt.Errorf("expected to received ping within 10 seconds")
	// 	}
	// }()

	if err := <-testError; err != nil {
		t.Fatal(err)
	}
}

type mockRealTimeChatWebsocketServer struct {
	state     State
	settings  Settings
	testError chan error
}

func (m *mockRealTimeChatWebsocketServer) createDefaultServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			m.testError <- err
			return
		}

		if r.Header.Get("Authorization") == m.settings.ValidOAuthToken {
			m.settings.Authorized = true
		}

		go func() {
			defer conn.Close()

			var (
				state         = ws.StateServerSide
				serverReader  = wsutil.NewReader(conn, state)
				serverWriter  = wsutil.NewWriter(conn, state, ws.OpText)
				serverEncoder = json.NewEncoder(serverWriter)
			)
			// Read forever until err
			for {
				header, err := serverReader.NextFrame()
				if err != nil {
					m.testError <- err
					return
				}

				lastWriteTime := time.Now()
				m.state.lastWriteFromClient = &lastWriteTime

				if !m.settings.Authorized {
					message := testWebsocketServerMessage{
						StatusCode: http.StatusUnauthorized,
						Message:    "Unable to verify the identity of the client",
					}

					if err := serverEncoder.Encode(message); err != nil {
						m.testError <- err
						return
					}

					if err := serverWriter.Flush(); err != nil {
						m.testError <- err
						return
					}

					serverWriter.Reset(conn, ws.StateServerSide, ws.OpClose)
					if err := serverEncoder.Encode(nil); err != nil {
						m.testError <- err
						return
					}

					if err := serverWriter.Flush(); err != nil {
						m.testError <- err
						return
					}

					if err := conn.Close(); err != nil {
						m.testError <- err
						return
					}

					return
				}

				if header.OpCode == ws.OpPing {
					serverWriter.Reset(conn, ws.StateServerSide, ws.OpPong)

					if err := serverEncoder.Encode(nil); err != nil {
						m.testError <- err
						return
					}

					if err := serverWriter.Flush(); err != nil {
						m.testError <- err
						return
					}
				}

				if err := serverReader.Discard(); err != nil {
					m.testError <- err
					return
				}

				m.testError <- nil
			}
		}()
	}))
}

type State struct {
	lastWriteFromClient *time.Time
}

type Settings struct {
	ValidOAuthToken string
	Authorized      bool
	WriteTimeout    time.Duration
}

type testWebsocketServerMessage struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
