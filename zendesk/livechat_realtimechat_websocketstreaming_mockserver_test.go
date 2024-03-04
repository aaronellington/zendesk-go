package zendesk_test

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type mockRealTimeChatWebsocketServer struct {
	state        state
	settings     settings
	conn         net.Conn
	handlers     frameHandlers
	queuedFrames []queuedFrame
	testError    chan error
}

func newMockRealTimeChatWebsocketServer(
	testErrorChan chan error,
	customState *state,
	customSettings *settings,
	customHandlers *frameHandlers,
) *mockRealTimeChatWebsocketServer {
	defaultState := state{
		successfulConnection: make(chan bool),
	}

	if customState != nil {
		defaultState = *customState
	}

	defaultSettings := settings{
		ValidOAuthToken: "Bearer fake-token",
	}

	if customSettings != nil {
		defaultSettings = *customSettings
	}

	return &mockRealTimeChatWebsocketServer{
		state:     defaultState,
		settings:  defaultSettings,
		testError: testErrorChan,
	}
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

			if err := m.read(); err != nil {
				mockServerError <- err
				return
			}
		}()

		if err := <-mockServerError; err != nil {
			m.testError <- err
			return
		}
	}))
}

func (m *mockRealTimeChatWebsocketServer) handleAuth() error {
	if m.settings.Authorized {
		m.state.successfulConnection <- m.settings.Authorized
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

		if header.OpCode.IsControl() {
			if err := m.handleControlFrame(payload, header.OpCode); err != nil {
				return err
			}
		}

		if header.OpCode.IsData() {
			if err := m.handleDataFrame(payload); err != nil {
				return err
			}
		}

		continue
	}
}

func (s *mockRealTimeChatWebsocketServer) handleDataFrame(
	data []byte,
) error {
	log.Println(string(data))

	return nil
}

func (m *mockRealTimeChatWebsocketServer) handleControlFrame(
	payload []byte,
	opCode ws.OpCode,
) error {
	receivedTime := time.Now()
	m.state.lastWriteFromClient = &receivedTime

	switch opCode {
	case ws.OpClose:
		return io.EOF
	case ws.OpPing:
		return m.handlers.ping(payload, ws.OpPong) // Something like this?
		// return m.write(payload, ws.OpPong)
	case ws.OpPong:
		return nil
	}

	return nil
}

type state struct {
	successfulConnection chan bool
	lastWriteFromClient  *time.Time
}

type settings struct {
	ValidOAuthToken string
	Authorized      bool
	WriteTimeout    time.Duration
}

type testWebsocketServerMessage struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type queuedFrame struct {
	payload []byte
	opCode  ws.OpCode
	delay   time.Duration
}

type frameHandlers struct {
	ping frameHandler
	pong frameHandler
}

type frameHandler func(payload []byte, opCode ws.OpCode) error
