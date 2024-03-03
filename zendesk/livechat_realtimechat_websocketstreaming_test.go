package zendesk_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
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

	timeout := time.NewTimer(time.Second * 20)
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a connection within timeout")
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

	timeout := time.NewTimer(time.Second * 200)
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a connection within timeout")
		// case successfulConnection := <-mockableZendeskRTCWebsocketServer.state.successfulConnection:
		// 	log.Println("successfulcon check")
		// 	if successfulConnection {
		// 		t.Fatal("expected to fail connection")
		// 	}

		// 	return
	}
}

type mockRealTimeChatWebsocketServer struct {
	state     State
	settings  Settings
	conn      net.Conn
	testError chan error
}

func (m *mockRealTimeChatWebsocketServer) createDefaultServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockServerError := make(chan error)

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			mockServerError <- err
			return
		}

		if r.Header.Get("Authorization") == m.settings.ValidOAuthToken {
			m.settings.Authorized = true
		}

		m.conn = conn

		go func() {
			if err := m.handleAuth(); err != nil {
				mockServerError <- err
				return
			}

			// if m.conn

			// if err := m.read(); err != nil {
			// 	mockServerError <- err
			// 	return
			// }
		}()

		if err := <-mockServerError; err != nil {
			m.testError <- err
			return
		}
	}))
}

func (m *mockRealTimeChatWebsocketServer) handleAuth() error {
	if m.settings.Authorized {
		// m.state.successfulConnection <- m.settings.Authorized
		return nil
	}

	message := testWebsocketServerMessage{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unable to verify the identity of the client",
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := m.write(messageBytes, ws.OpText); err != nil {
		return err
	}

	if err := m.write(nil, ws.OpClose); err != nil {
		return err
	}

	// m.state.successfulConnection <- m.settings.Authorized

	if err := m.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (m *mockRealTimeChatWebsocketServer) write(payload []byte, opCode ws.OpCode) error {
	writer := wsutil.NewWriter(m.conn, ws.StateServerSide, opCode)
	_, err := writer.Write(payload)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func (m *mockRealTimeChatWebsocketServer) read() error {
	for {
		header, err := ws.ReadHeader(m.conn)
		if err != nil {
			return err
		}

		lastWriteFromClient := time.Now()
		m.state.lastWriteFromClient = &lastWriteFromClient

		payload := make([]byte, header.Length)
		_, err = io.ReadFull(m.conn, payload)
		if err != nil {
			return err
		}

		if header.Masked {
			ws.Cipher(payload, header.Mask, 0)
		}

		// if header.OpCode.IsControl() {
		// 	if err := s.handleControlFrame(payload, header.OpCode); err != nil {
		// 		return err
		// 	}
		// }

		// if header.OpCode.IsData() {
		// 	if err := s.handleDataFrame(payload); err != nil {
		// 		return err
		// 	}
		// }

		continue
	}
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
