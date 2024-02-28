package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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
	client    *client   //
	rtcWSHost string    // The host for the LiveChat RealTimeChat Websocket server - present here so that it can be overridden in tests
	conn      *net.Conn //
}

type wsChatCache struct {
	individualDepartments *utils.MemoryCacheInstance[GroupID, WebsocketChatMetricData]
}

type wsAgentCache struct {
	individualDepartments *utils.MemoryCacheInstance[UserID, WebsocketAgentMetricData]
}
type wsConnMetadata struct {
	mutex           *sync.Mutex
	connStarted     *time.Time
	sentPing        *time.Time
	sentData        *time.Time
	receivedControl *time.Time
	receivedData    *time.Time
}

func getMostRecentTime(t1, t2 *time.Time) *time.Duration {
	if t1 == nil && t2 == nil {
		return nil
	}

	var t2d time.Duration
	if t1 == nil {
		t2d = time.Since(*t2)
		log.Print("return t2")
		return &t2d
	}

	t1d := time.Since(*t1)
	if t1d < t2d {
		log.Print("return t2 t1 is less")
		return &t2d
	}
	log.Print("return t1")
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
	if s.wsClient.conn != nil {
		return nil
	}

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

	s.wsClient.conn = &conn

	connectionStartedTime := time.Now()

	s.wsCache.metadata.connStarted = &connectionStartedTime

	return nil
}

func (s *RealTimeChatStreamingService) retryableErr() bool {
	return true
}

func (s *RealTimeChatStreamingService) cleanConnection() {
	s.wsClient.conn = nil

	newCache := &wsCache{
		chat: &wsChatCache{
			individualDepartments: &utils.MemoryCacheInstance[GroupID, WebsocketChatMetricData]{},
		},
		agent: &wsAgentCache{
			individualDepartments: &utils.MemoryCacheInstance[UserID, WebsocketAgentMetricData]{},
		},
		metadata: &wsConnMetadata{
			mutex: &sync.Mutex{},
		},
	}

	s.wsCache = newCache
}

func (s *RealTimeChatStreamingService) ConnectToWebsocket(parentCtx context.Context) error {
	ctx, cancelHandler := context.WithCancel(parentCtx)
	defer func() {
		s.cleanConnection()
		cancelHandler()
	}()

	if err := s.initiateWebsocketConnection(ctx); err != nil {
		return err
	}

	errorChan := make(chan error)
	go func() {
		if err := s.ping(ctx); err != nil {
			errorChan <- err
		}
	}()

	go func() {
		if err := s.read(ctx); err != nil {
			errorChan <- err
		}
	}()

	return <-errorChan
}

func (s *RealTimeChatStreamingService) ping(ctx context.Context) error {
	if err := s.write(ctx, nil, ws.OpPing); err != nil {
		return err
	}

	firstPing := time.Now()
	s.wsCache.metadata.sentPing = &firstPing

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case currentPingSentTime := <-ticker.C:
			if err := s.write(ctx, nil, ws.OpPing); err != nil {
				return err
			}

			s.wsCache.metadata.sentPing = &currentPingSentTime

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *RealTimeChatStreamingService) read(ctx context.Context) error {
	if s.wsClient.conn == nil {
		return ErrRealTimeChatWebsocketConnectionIsNil
	}

	reader := wsutil.NewClientSideReader(*s.wsClient.conn)
	decoder := json.NewDecoder(reader)
	for {

		header, err := reader.NextFrame()
		if err != nil {
			return err
		}

		if header.OpCode.IsControl() {

			if err := s.handleControlFrame(reader, header); err != nil {
				return err
			}

			continue
		}

		if header.OpCode.IsData() {
			b := map[string]any{}
			if err := decoder.Decode(&b); err != nil {
				return err
			}

			messageBytes, err := json.Marshal(b)
			if err != nil {
				return err
			}

			if err := s.handleWebsocketMessage(ctx, messageBytes); err != nil {
				return err
			}

			continue
		}
	}
}

func (s *RealTimeChatStreamingService) handleWebsocketMessage(
	ctx context.Context,
	messageBytes []byte,
) error {
	placeholder := RealTimeChatStreamingResponse{}
	if err := json.Unmarshal(messageBytes, &placeholder); err != nil {
		return err
	}

	status, ok := placeholder["status_code"]
	if !ok {
		return ErrRealTimeChatWebsocketUnsupportedIncomingMessage
	}

	if status == float64(401) {
		return ErrRealTimeChatWebsocketUnauthenticated
	}

	return nil
}

func (s *RealTimeChatStreamingService) write(ctx context.Context, payload any, opCode ws.OpCode) error {
	if s.wsClient.conn == nil {
		return ErrRealTimeChatWebsocketConnectionIsNil
	}

	writer := wsutil.NewWriter(*s.wsClient.conn, ws.StateClientSide, ws.OpText)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(payload); err != nil {
		return err
	}

	return writer.Flush()
}

func (s *RealTimeChatStreamingService) handleControlFrame(reader *wsutil.Reader, header ws.Header) error {
	s.wsCache.metadata.mutex.Lock()
	defer s.wsCache.metadata.mutex.Unlock()

	receivedTime := time.Now()
	s.wsCache.metadata.receivedControl = &receivedTime

	switch header.OpCode {
	case ws.OpClose:
		return io.EOF
	case ws.OpPong:
		if err := reader.Discard(); err != nil {
			return err
		}
	}

	return nil
}

func (s *RealTimeChatStreamingService) WebsocketReady() bool {
	return s.wsClient.conn != nil
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

	if err := s.write(ctx, payload, ws.OpText); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	return nil
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

type RealTimeChatStreamingResponse map[string]any
