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
		state: State{
			successfulConnection: make(chan bool),
		},
		settings: Settings{
			ValidOAuthToken: "Bearer fake-token",
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

	timeout := time.NewTimer(time.Second * 2)
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a successful connection within timeout")
	case successfulConnection := <-mockableZendeskRTCWebsocketServer.state.successfulConnection:
		if !successfulConnection {
			t.Fatal("did not connect successfully")
		}

		return
	}
}

func TestRealTimeChatWebsocketStreaming_Connect_401(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	// This is our test mockserver
	mockableZendeskRTCWebsocketServer := &mockRealTimeChatWebsocketServer{
		state: State{
			successfulConnection: make(chan bool),
		},
		settings: Settings{
			ValidOAuthToken: "No valid token",
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

	timeout := time.NewTimer(time.Second * 2)
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a successful connection within timeout")
	case successfulConnection := <-mockableZendeskRTCWebsocketServer.state.successfulConnection:
		if successfulConnection {
			t.Fatal("expected to fail connection")
		}

		return
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

			m.state.successfulConnection <- m.settings.Authorized

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

					_, err := serverWriter.Write([]byte{})
					if err != nil {
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

				m.testError <- err
			}
		}()
	}))
}

type State struct {
	successfulConnection chan bool
	lastWriteFromClient  *time.Time
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
