package zendesk_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/equalsgibson/concur/concur"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type mockRealTimeChatWebsocketServer struct {
	state     state
	settings  settings
	conn      net.Conn
	connError chan error
	handlers  frameHandlers
	history   struct {
		receivedFrames []frame
		sentFrames     []frame
	}
	queuedFrames struct {
		subscribe   map[string][]queuedFrame
		unsubscribe map[string][]queuedFrame
	}
	testError chan error
}

func newMockRealTimeChatWebsocketServer(
	t *testing.T,
	customSettings settings,
) (*mockRealTimeChatWebsocketServer, string) {
	mockserver := &mockRealTimeChatWebsocketServer{
		state: state{
			ValidOAuthToken: "Bearer fake-token",
		},
		settings:  customSettings,
		testError: make(chan error),
		connError: make(chan error),
	}

	return mockserver, mockserver.createDefaultServer(t)
}

func (m *mockRealTimeChatWebsocketServer) createDefaultServer(t *testing.T) string {
	svr := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockServerError := make(chan error)

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			mockServerError <- err
			return
		}

		if r.Header.Get("Authorization") == m.state.ValidOAuthToken {
			m.state.Authorized = true
		}

		m.conn = conn

		log.Println("Connected to Server!!")

		if err := m.handleAuth(); err != nil {
			log.Printf("Unauthenticated: %s", err)
			m.conn.Close()
			return
		}

		serverReader := concur.NewAsyncReader(func(ctx context.Context) ([]byte, error) {
			return m.read()
		})

		go serverReader.Loop(r.Context())
		defer serverReader.Close()

		for {
			select {
			case update := <-serverReader.Updates():
				if update.Err != nil {
					log.Printf("Update error: %s", err)
					m.conn.Close()
					return
				}

				if err := s.handleFrame(update.Item); err != nil {
					log.Printf("Handle update error: %s", err)
					m.conn.Close()
					return
				}
			}
		}

		// 	if header.OpCode.IsControl() {
		// 		if err := m.handleControlFrame(payload, header.OpCode); err != nil {
		// 			return nil, err
		// 		}
		// 	}

		// 	if header.OpCode.IsData() {
		// 		if err := m.handleDataFrame(payload); err != nil {
		// 			return nil, err
		// 		}
		// 	}
		// }

		// if header.OpCode.IsControl() {
		// 	if err := m.handleControlFrame(payload, header.OpCode); err != nil {
		// 		return nil, err
		// 	}
		// }

		// if header.OpCode.IsData() {
		// 	if err := m.handleDataFrame(payload); err != nil {
		// 		return nil, err
		// 	}
		// }

		// go func() {

		// 	// for _, frame := range m.queuedFrames {
		// 	// 	time.Sleep(frame.delay)
		// 	// 	if err := m.write(frame.payload, frame.opCode); err != nil {
		// 	// 		mockServerError <- err
		// 	// 		return
		// 	// 	}
		// 	// }
		// }()

		if err := <-mockServerError; err != nil {
			log.Printf("Closing connection because of read/write error: %s", err)
			m.conn.Close()
			return
		}
	}))

	svr.Start()
	t.Cleanup(
		svr.Close,
	)

	return svr.URL
}

func (s *RealTimeChatStreamingService) handleFrame(frame realTimeChatStreamingFrame) error {
	if frame.opCode.IsControl() {
		if err := s.handleControlFrame(frame); err != nil {
			return err
		}
	}

	if frame.opCode.IsData() {
		if err := s.handleDataFrame(frame); err != nil {
			return err
		}
	}

	return nil
}

func (m *mockRealTimeChatWebsocketServer) handleAuth() error {
	if m.state.Authorized {
		message := websocketServerMessage{
			StatusCode: http.StatusOK,
			Message: struct {
				Authenticated bool `json:"authenticated"`
			}{
				Authenticated: true,
			},
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			return err
		}

		if err := m.write(messageBytes, ws.OpText); err != nil {
			return err
		}

		return nil
	}

	message := websocketServerMessage{
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

	log.Println("Closing inside the handleAuth func!")

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

	m.history.sentFrames = append(m.history.sentFrames, frame{
		payload: payload,
		opCode:  opCode,
		time:    time.Now(),
	})

	return writer.Flush()
}

func (m *mockRealTimeChatWebsocketServer) read() ([]byte, error) {
	header, err := ws.ReadHeader(m.conn)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(m.conn, payload)
	if err != nil {
		return nil, err
	}

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	m.history.receivedFrames = append(m.history.receivedFrames, frame{
		payload: payload,
		opCode:  header.OpCode,
		time:    time.Now(),
	})

	return payload, nil
}

func (s *mockRealTimeChatWebsocketServer) handleDataFrame(
	data []byte,
) error {
	target := zendesk.Subscription{}
	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	if target.Action == "subscribe" {
		go func() {
			queuedFrames, ok := s.queuedFrames.subscribe[target.Topic]
			if !ok {
				s.connError <- errors.New("No queued frames to be sent for topic!")
				return
			}

			for _, frame := range queuedFrames {
				time.Sleep(frame.delay)

				if err := s.write(frame.payload, frame.opCode); err != nil {
					s.connError <- err
					return
				}
			}
		}()
	}

	return nil
}

func (m *mockRealTimeChatWebsocketServer) handleControlFrame(
	payload []byte,
	opCode ws.OpCode,
) error {
	switch opCode {
	case ws.OpClose:
		return io.EOF
	case ws.OpPing:
		receivedTime := time.Now()
		m.state.LastPingReceived = &receivedTime

		return m.write(payload, ws.OpPong)
	case ws.OpPong:
		return nil
	}

	return nil
}

type state struct {
	Authorized       bool
	ValidOAuthToken  string
	LastPingReceived *time.Time
}

type settings struct {
	ShouldResponseToPing bool
}

type websocketServerMessage struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

type queuedFrame struct {
	payload []byte
	opCode  ws.OpCode
	delay   time.Duration
}

type realTimeChatStreamingFrameHistory struct {
	realTimeChatStreamingFrame
	time time.Time
}

type realTimeChatStreamingFrame struct {
	opCode  ws.OpCode
	payload []byte
}

type frameHandlers struct {
	ping frameHandler
	pong frameHandler
}

type frameHandler func(payload []byte, opCode ws.OpCode) error
