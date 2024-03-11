package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/utils"
	"github.com/equalsgibson/concur/concur"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WebsocketCloseCode int

const (
	WebsocketCloseCodeNormal          WebsocketCloseCode = 1000
	WebsocketCloseCodeGoingAway       WebsocketCloseCode = 1001
	WebsocketCloseCodeProtocolErr     WebsocketCloseCode = 1002
	WebsocketCloseCodeUnsupportedData WebsocketCloseCode = 1003
	WebsocketCloseCodeNoCodeReceived  WebsocketCloseCode = 1005
)

const RealTimeChatStreamingHost string = "wss://rtm.zopim.com/stream"

var (
	ErrRealTimeChatWebsocketUnauthenticated            error = errors.New("unauthenticated")
	ErrRealTimeChatWebsocketUnsupportedIncomingMessage error = errors.New("unsupported incoming message")
	ErrRealTimeChatWebsocketConnectionIsNil            error = errors.New("realtimechat websocket connection is nil")
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
type RealTimeChatStreamingService struct {
	wsClient *wsClient
	wsCache  *wsCache
}

type wsCache struct {
	chat     *wsChatCache
	agent    *wsAgentCache
	metadata *wsConnMetadata
}

type wsClient struct {
	client          *client //
	conn            net.Conn
	connEstablished chan bool
	connMutex       *sync.Mutex
	rtcWSHost       string // The host for the LiveChat RealTimeChat Websocket server - present here so that it can be overridden in tests
}

type wsChatCache struct {
	individualDepartments *utils.MemoryCacheInstance[GroupID, WebsocketChatMetricData]
}

type wsAgentCache struct {
	individualDepartments *utils.MemoryCacheInstance[UserID, WebsocketAgentMetricData]
}

type wsConnMetadata struct {
	mutex             *sync.Mutex
	connStarted       *time.Time
	connAuthenticated *time.Time
	sentPing          *time.Time
	sentData          *time.Time
	receivedControl   *time.Time
	receivedData      *time.Time
}

func getMostRecentTime(t1, t2 *time.Time) *time.Duration {
	if t1 == nil && t2 == nil {
		return nil
	}

	var t2d time.Duration
	if t1 == nil {
		t2d = time.Since(*t2)
		return &t2d
	}

	t1d := time.Since(*t1)
	if t1d < t2d {
		return &t2d
	}

	return &t1d
}

type WebsocketChatMetricData struct {
	MissedChats        *ChatMetricWindow               `json:"missed_chats"`
	ChatDurationMax    *uint64                         `json:"chat_duration_max"`
	SatisfactionBad    *ChatMetricWindow               `json:"satisfaction_bad"`
	ActiveChats        uint64                          `json:"active_chats"`
	SatisfactionGood   *ChatMetricWindow               `json:"satisfaction_good"`
	IncomingChats      uint64                          `json:"incoming_chats"`
	AssignedChats      uint64                          `json:"assigned_chats"`
	ChatDurationAvg    *uint64                         `json:"chat_duration_avg"`
	WaitingTimeAvg     *uint64                         `json:"waiting_time_avg"`
	ResponseTimeAvg    *uint64                         `json:"response_time_avg"`
	WaitingTimeMax     *uint64                         `json:"waiting_time_max"`
	ResponseTimeMax    *uint64                         `json:"response_time_max"`
	Subscriptions      WebsocketChatMetricSubscription `json:"subscriptions"`
	LastUpdateReceived *time.Time                      `json:"last_update_received"`
}

type WebsocketChatMetricSubscription struct {
	MissedChats30      bool `json:"missed_chats_30"`
	MissedChats60      bool `json:"missed_chats_60"`
	ChatDurationMax    bool `json:"chat_duration_max"`
	SatisfactionBad30  bool `json:"satisfaction_bad_30"`
	SatisfactionBad60  bool `json:"satisfaction_bad_60"`
	ActiveChats        bool `json:"active_chats"`
	SatisfactionGood30 bool `json:"satisfaction_good_30"`
	SatisfactionGood60 bool `json:"satisfaction_good_60"`
	IncomingChats      bool `json:"incoming_chats"`
	AssignedChats      bool `json:"assigned_chats"`
	ChatDurationAvg    bool `json:"chat_duration_avg"`
	WaitingTimeAvg     bool `json:"waiting_time_avg"`
	ResponseTimeAvg    bool `json:"response_time_avg"`
	WaitingTimeMax     bool `json:"waiting_time_max"`
	ResponseTimeMax    bool `json:"response_time_max"`
}

type WebsocketAgentMetricData struct {
	AgentsOnline       *uint64                          `json:"agents_online"`
	AgentsAway         *uint64                          `json:"agents_away"`
	AgentsInvisible    *uint64                          `json:"agents_invisible"`
	Subscriptions      WebsocketAgentMetricSubscription `json:"subscriptions"`
	LastUpdateReceived *time.Time                       `json:"last_update_received"`
}

type WebsocketAgentMetricSubscription struct {
	AgentsOnline    bool `json:"agents_online"`
	AgentsAway      bool `json:"agents_away"`
	AgentsInvisible bool `json:"agents_invisible"`
}

func (s *RealTimeChatStreamingService) initiateWebsocketConnection(ctx context.Context) error {
	if err := s.wsClient.client.GetAccessToken(ctx); err != nil {
		return err
	}

	headers := ws.HandshakeHeaderHTTP{}
	headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", s.wsClient.client.chatToken.AccessToken)}

	dialer := ws.Dialer{
		Header: headers,
	}

	conn, _, _, err := dialer.Dial(ctx, s.wsClient.rtcWSHost)
	if err != nil {
		return err
	}

	connectionStartedTime := time.Now()
	s.wsCache.metadata.connStarted = &connectionStartedTime

	s.wsClient.conn = conn
	// s.wsClient.connEstablished <- true

	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	conn.Close()
	// 	s.wsClient.conn.Close()
	// }()

	return nil
}

func (s *RealTimeChatStreamingService) ConnectToWebsocket(parentCtx context.Context) error {
	ctx, _ := context.WithCancelCause(parentCtx)

	if err := s.initiateWebsocketConnection(ctx); err != nil {
		return err
	}
	defer s.wsClient.conn.Close()

	reader := concur.NewAsyncReader[[]byte](
		func(ctx context.Context) ([]byte, error) {
			return s.read()
		},
	)

	keepalive := time.NewTicker(time.Second)
	defer keepalive.Stop()
	go func() {
		for range keepalive.C {
			if err := s.ping(); err != nil {
				reader.Close()
				return
			}
		}
	}()

	// pongMonitor := time.NewTicker(time.Second)
	// defer pongMonitor.Stop()
	// go func() {
	// 	for range pongMonitor.C {
	// 		s.ping()
	// 	}
	// }()

	go reader.Loop(ctx)

	var err error
	for update := range reader.Updates() {
		if update.Err != nil {
			err = update.Err
			reader.Close()
		}

		// if header.OpCode.IsControl() {
		// 	if err := s.handleControlFrame(payload, header.OpCode); err != nil {
		// 		return nil, err
		// 	}
		// }

		// if header.OpCode.IsData() {
		// 	if err := s.handleDataFrame(payload); err != nil {
		// 		return nil, err
		// 	}
		// }

	}

	return err
}

func (s *RealTimeChatStreamingService) CloseConnection(closeStatusCode WebsocketCloseCode, closeReason string) error {
	closeFrame := RealTimeChatStreamingCloseFrame{
		Code:   closeStatusCode,
		Reason: closeReason,
	}

	frameBytes, err := json.Marshal(closeFrame)
	if err != nil {
		return err
	}

	if err := s.write(frameBytes, ws.OpClose); err != nil {
		return err
	}

	return nil
}

func (s *RealTimeChatStreamingService) ping() error {
	if err := s.write(nil, ws.OpPing); err != nil {
		return err
	}

	firstPing := time.Now()
	s.wsCache.metadata.sentPing = &firstPing

	return nil
}

func (s *RealTimeChatStreamingService) read() ([]byte, error) {
	header, err := ws.ReadHeader(s.wsClient.conn)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(s.wsClient.conn, payload)
	if err != nil {
		return nil, err
	}

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	return payload, nil

	// return payload
}

func (s *RealTimeChatStreamingService) write(payload []byte, opCode ws.OpCode) error {
	if err := s.connectionEstablished(); err != nil {
		return err
	}

	writer := wsutil.NewWriter(s.wsClient.conn, ws.StateClientSide, opCode)
	_, err := writer.Write(payload)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func (s *RealTimeChatStreamingService) connectionEstablished() error {
	if s.wsClient.conn != nil {
		return nil
	}

	s.wsClient.connMutex.Lock()
	defer s.wsClient.connMutex.Unlock()

	if s.wsClient.conn != nil {
		return nil
	}

	timeout := time.NewTimer(time.Second * 15)
	for {
		select {
		case <-timeout.C:
			return fmt.Errorf("timeout reached")
		case <-s.wsClient.connEstablished:
			return nil
		}
	}
}

func (s *RealTimeChatStreamingService) handleControlFrame(
	payload []byte,
	opCode ws.OpCode,
) error {
	receivedTime := time.Now()
	s.wsCache.metadata.receivedControl = &receivedTime

	switch opCode {
	case ws.OpPing:
		return s.write(payload, ws.OpPong)
	case ws.OpPong:
		return nil
	case ws.OpClose:
		if err := s.write(payload, ws.OpClose); err != nil {
			return err
		}

		return io.EOF
	}

	return nil
}

func (s *RealTimeChatStreamingService) handleDataFrame(
	data []byte,
) error {
	fmt.Println("From Server: ", string(data))

	type status struct {
		StatusCode int `json:"status_code"`
	}

	frame := status{}
	if err := json.Unmarshal(data, &frame); err != nil {
		return err
	}

	switch frame.StatusCode {
	case http.StatusUnauthorized:
		return ErrRealTimeChatWebsocketUnauthenticated
	case http.StatusInternalServerError:
		return errors.New("error, unsupported message sent to server: " + string(data))
	case http.StatusOK:

	}

	return nil
}

type Subscription struct {
	Topic  string `json:"topic"`
	Action string `json:"action"`
}

func (s *RealTimeChatStreamingService) SubscribeToAgentMetric(ctx context.Context, metric LiveChatMetricKeyAgent) error {
	payload := Subscription{
		Topic:  fmt.Sprintf("agents.%s", metric),
		Action: "subscribe",
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(bytes, ws.OpText); err != nil {
		return err
	}

	fmt.Println(string(bytes))

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	return nil
}

func (s *RealTimeChatStreamingService) GetAgentMetric(ctx context.Context) (map[UserID]WebsocketAgentMetricData, error) {
	items, err := s.wsCache.agent.individualDepartments.GetAll()
	if err != nil {
		return nil, err
	}

	return items.Items, nil
}

func (s *RealTimeChatStreamingService) GetTimeSinceLastFrameReceived() *time.Duration {
	return getMostRecentTime(s.wsCache.metadata.receivedControl, s.wsCache.metadata.receivedData)
}

func (s *RealTimeChatStreamingService) GetTimeSinceLastFrameSent() *time.Duration {
	return getMostRecentTime(s.wsCache.metadata.sentPing, s.wsCache.metadata.sentData)
}

type RealTimeChatStreamingMessage struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type RealTimeChatStreamingContent[payload any] struct {
	StatusCode int     `json:"status_code"`
	Content    payload `json:"content"`
}

type RealTimeChatStreamingContentAgentMetric struct {
	Topic        string                          `json:"topic"`
	Data         map[LiveChatMetricKeyAgent]uint `json:"data"`
	Type         string                          `json:"type"`
	DepartmentID *GroupID                        `json:"department_id"`
}

type RealTimeChatStreamingCloseFrame struct {
	Code   WebsocketCloseCode `json:"code"`
	Reason string             `json:"reason,omitempty"`
}
