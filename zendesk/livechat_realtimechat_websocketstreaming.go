package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
type RealTimeChatStreamingService struct {
	client         *client
	wsConn         *net.Conn
	wsChatCache    *wsChatCache
	wsAgentCache   *wsAgentCache
	wsConnMetadata *wsConnMetadata
}

type wsChatCache struct {
	individualDepartments *utils.MemoryCacheInstance[GroupID, WebsocketChatMetricData]
	globalMetrics         *WebsocketChatMetricData
}

type wsAgentCache struct {
	individualDepartments *utils.MemoryCacheInstance[UserID, WebsocketAgentMetricData]
	globalMetrics         *WebsocketAgentMetricData
}
type wsConnMetadata struct {
	mutex           *sync.Mutex
	connStarted     *time.Time
	sentPing        *time.Time
	receivedControl *time.Time
	receivedData    *time.Time
}

func (w *wsConnMetadata) getTimeSinceMostRecentFrameReceived() *time.Duration {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	var controlFrameTime time.Duration
	var dataFrameTime time.Duration

	if w.receivedControl != nil {
		controlFrameTime = time.Since(*w.receivedControl)
	}

	if w.receivedData != nil {
		dataFrameTime = time.Since(*w.receivedData)
	}

	if controlFrameTime == 0 && dataFrameTime == 0 {
		return nil
	}

	if controlFrameTime == 0 {
		return &dataFrameTime
	}

	return &controlFrameTime
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

func (s *RealTimeChatStreamingService) InitiateWebsocketConnection(ctx context.Context) error {
	if s.wsConn != nil {
		return nil
	}

	if err := s.client.GetAccessToken(ctx); err != nil {
		return err
	}

	headers := ws.HandshakeHeaderHTTP{}
	headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", s.client.chatToken.AccessToken)}

	dialer := ws.Dialer{
		Header: headers,
	}

	conn, _, _, err := dialer.Dial(ctx, "wss://rtm.zopim.com/stream")
	if err != nil {
		return err
	}

	s.wsConn = &conn

	s.wsConnMetadata.mutex.Lock()
	defer s.wsConnMetadata.mutex.Unlock()

	connectionStartedTime := time.Now()

	s.wsConnMetadata.connStarted = &connectionStartedTime

	return nil
}

func (s *RealTimeChatStreamingService) ConnectToWebsocket(parentCtx context.Context) error {
	ctx, cancelHandler := context.WithCancel(parentCtx)
	defer func() {
		cancelHandler()
	}()

	if err := s.InitiateWebsocketConnection(ctx); err != nil {
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

	// go func() {
	// 	if err := s.write(ctx); err != nil {
	// 		errorChan <- err
	// 	}
	// }()

	return <-errorChan
}

func (s *RealTimeChatStreamingService) ping(ctx context.Context) error {
	writer := wsutil.NewWriter(*s.wsConn, ws.StateClientSide, ws.OpPing)
	encoder := json.NewEncoder(writer)

	// Send first ping immediately
	if err := encoder.Encode(nil); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	firstPing := time.Now()
	s.wsConnMetadata.sentPing = &firstPing

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case currentPingSentTime := <-ticker.C:
			if err := encoder.Encode(nil); err != nil {
				return err
			}

			if err := writer.Flush(); err != nil {
				return err
			}

			s.wsConnMetadata.sentPing = &currentPingSentTime

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *RealTimeChatStreamingService) read(ctx context.Context) error {
	reader := wsutil.NewClientSideReader(*s.wsConn)
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

			fmt.Println(b)

			continue
		}
	}
}

func (s *RealTimeChatStreamingService) write(ctx context.Context) error {
	writer := wsutil.NewWriter(*s.wsConn, ws.StateClientSide, ws.OpText)
	encoder := json.NewEncoder(writer)
	for {

		// b := Subscription{
		// 	Topic:  "chats.incoming_chats",
		// 	Action: "subscribe",
		// }

		if err := encoder.Encode(nil); err != nil {
			return err
		}

		if err := writer.Flush(); err != nil {
			return err
		}

		time.Sleep(time.Second * 5)
	}
}

func (s *RealTimeChatStreamingService) handleControlFrame(reader *wsutil.Reader, header ws.Header) error {
	s.wsConnMetadata.mutex.Lock()
	defer s.wsConnMetadata.mutex.Unlock()

	receivedTime := time.Now()
	s.wsConnMetadata.receivedControl = &receivedTime

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

type Subscription struct {
	Topic  string `json:"topic"`
	Action string `json:"action"`
}

func (s *RealTimeChatStreamingService) SubscribeToAgentMetric()                     {}
func (s *RealTimeChatStreamingService) SubscribeToChatMetric()                      {}
func (s *RealTimeChatStreamingService) SubscribeToChatMetricForSpecificTimeWindow() {}

func (s *RealTimeChatStreamingService) GetTimeSinceLastFrameReceived() *time.Duration {
	return s.wsConnMetadata.getTimeSinceMostRecentFrameReceived()
}

func (s *RealTimeChatStreamingService) GetTimeSinceLastFrameSent() *time.Duration {
	return s.wsConnMetadata.getTimeSinceMostRecentFrameReceived()
}
