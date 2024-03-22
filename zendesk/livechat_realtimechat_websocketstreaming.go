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
	ErrRealTimeChatWebsocketUnsupportedIncomingMessage error = errors.New("unsupported message received by server")
	ErrRealTimeChatWebsocketUnsupportedSentMessage     error = errors.New("unsupported message sent to server")
	ErrRealTimeChatWebsocketConnectionIsNil            error = errors.New("realtimechat websocket connection is nil")
)

const realTimeChatStreamingGlobalDepartmentID GroupID = 0

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
type RealTimeChatStreamingService struct {
	wsClient *wsClient
	wsCache  *wsCache
}

type wsCache struct {
	chat     *utils.MemoryCacheInstance[GroupID, WebsocketChatMetricData]
	agent    *utils.MemoryCacheInstance[GroupID, WebsocketAgentMetricData]
	metadata *wsConnMetadata
}

type wsClient struct {
	client    *client
	conn      net.Conn
	connMutex *sync.Mutex
	rtcWSHost string // The host for the LiveChat RealTimeChat Websocket server - present here so that it can be overridden in tests
}

type wsConnMetadata struct {
	mutex           *sync.Mutex
	connStarted     *time.Time
	sentPing        *time.Time
	sentData        *time.Time
	receivedControl *time.Time
	receivedData    *time.Time
}

type WebsocketChatMetricData struct {
	MissedChats        ChatMetricWindow                `json:"missed_chats"`
	ChatDurationMax    *uint64                         `json:"chat_duration_max"`
	SatisfactionBad    ChatMetricWindow                `json:"satisfaction_bad"`
	ActiveChats        uint64                          `json:"active_chats"`
	SatisfactionGood   ChatMetricWindow                `json:"satisfaction_good"`
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

func (w *WebsocketChatMetricData) patchData(update map[LiveChatMetricKeyChat]any) error {
	updated := time.Now()

	for dataPointKey, rawData := range update {
		if LiveChatMetricKeyChatWindow(dataPointKey) == LiveChatMetricKeyMissedChats ||
			LiveChatMetricKeyChatWindow(dataPointKey) == LiveChatMetricKeySatisfactionBad ||
			LiveChatMetricKeyChatWindow(dataPointKey) == LiveChatMetricKeySatisfactionGood {

			windowMetric := map[string]int64{}
			rawDataBytes, err := json.Marshal(rawData)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(rawDataBytes, &windowMetric); err != nil {
				return err
			}

			switch LiveChatMetricKeyChatWindow(dataPointKey) {
			case LiveChatMetricKeyMissedChats:
				if data, ok := windowMetric["30"]; ok {
					w.MissedChats.ThirtyMinuteWindow = data
				}
				if data, ok := windowMetric["60"]; ok {
					w.MissedChats.SixtyMinuteWindow = data
				}

			case LiveChatMetricKeySatisfactionBad:
				if data, ok := windowMetric["30"]; ok {
					w.SatisfactionBad.ThirtyMinuteWindow = data
				}
				if data, ok := windowMetric["60"]; ok {
					w.SatisfactionBad.SixtyMinuteWindow = data
				}

			case LiveChatMetricKeySatisfactionGood:
				if data, ok := windowMetric["30"]; ok {
					w.SatisfactionGood.ThirtyMinuteWindow = data
				}
				if data, ok := windowMetric["60"]; ok {
					w.SatisfactionGood.SixtyMinuteWindow = data
				}
			}
		} else {
			// NOTE: Default for JSON is to return float64 for numbers.
			dataFloat, ok := rawData.(float64)
			if !ok {
				return errors.New("could not convert data to float64")
			}

			data := uint64(dataFloat)

			switch dataPointKey {
			case LiveChatMetricKeyChatDurationMax:
				w.ChatDurationMax = &data

			case LiveChatMetricKeyActiveChats:
				w.ActiveChats = data

			case LiveChatMetricKeyIncomingChats:
				w.IncomingChats = data

			case LiveChatMetricKeyAssignedChats:
				w.AssignedChats = data

			case LiveChatMetricKeyChatDurationAvg:
				w.ChatDurationAvg = &data

			case LiveChatMetricKeyWaitingTimeAvg:
				w.WaitingTimeAvg = &data

			case LiveChatMetricKeyResponseTimeAvg:
				w.ResponseTimeAvg = &data

			case LiveChatMetricKeyWaitingTimeMax:
				w.WaitingTimeMax = &data

			case LiveChatMetricKeyResponseTimeMax:
				w.ResponseTimeMax = &data
			}
		}

	}

	w.LastUpdateReceived = &updated

	return nil
}

// NOTE: map[string]bool is to account for ChatMetricWindow subs.
func (w *WebsocketChatMetricData) patchSubscription(update map[string]bool) {
	for dataPointKey, data := range update {
		switch dataPointKey {
		case fmt.Sprintf("%s30", LiveChatMetricKeyMissedChats):
			w.Subscriptions.MissedChats30 = data
		case fmt.Sprintf("%s60", LiveChatMetricKeyMissedChats):
			w.Subscriptions.MissedChats60 = data

		case string(LiveChatMetricKeyChatDurationMax):
			w.Subscriptions.ChatDurationMax = data

		case fmt.Sprintf("%s30", LiveChatMetricKeySatisfactionBad):
			w.Subscriptions.SatisfactionBad30 = data
		case fmt.Sprintf("%s60", LiveChatMetricKeySatisfactionBad):
			w.Subscriptions.SatisfactionBad60 = data

		case string(LiveChatMetricKeyActiveChats):
			w.Subscriptions.ActiveChats = data

		case fmt.Sprintf("%s30", LiveChatMetricKeySatisfactionGood):
			w.Subscriptions.SatisfactionGood30 = data
		case fmt.Sprintf("%s60", LiveChatMetricKeySatisfactionGood):
			w.Subscriptions.SatisfactionGood60 = data

		case string(LiveChatMetricKeyIncomingChats):
			w.Subscriptions.IncomingChats = data

		case string(LiveChatMetricKeyAssignedChats):
			w.Subscriptions.AssignedChats = data

		case string(LiveChatMetricKeyChatDurationAvg):
			w.Subscriptions.ChatDurationAvg = data

		case string(LiveChatMetricKeyWaitingTimeAvg):
			w.Subscriptions.WaitingTimeAvg = data

		case string(LiveChatMetricKeyResponseTimeAvg):
			w.Subscriptions.ResponseTimeAvg = data

		case string(LiveChatMetricKeyWaitingTimeMax):
			w.Subscriptions.WaitingTimeMax = data

		case string(LiveChatMetricKeyResponseTimeMax):
			w.Subscriptions.ResponseTimeMax = data
		}
	}
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
	AgentsOnline       uint64                           `json:"agents_online"`
	AgentsAway         uint64                           `json:"agents_away"`
	AgentsInvisible    uint64                           `json:"agents_invisible"`
	Subscriptions      WebsocketAgentMetricSubscription `json:"subscriptions"`
	LastUpdateReceived *time.Time                       `json:"last_update_received"`
}

type WebsocketAgentMetricSubscription struct {
	AgentsOnline    bool `json:"agents_online"`
	AgentsAway      bool `json:"agents_away"`
	AgentsInvisible bool `json:"agents_invisible"`
}

func (w *WebsocketAgentMetricData) patchData(update map[LiveChatMetricKeyAgent]uint64) {
	updated := time.Now()

	for dataPointKey, data := range update {
		switch dataPointKey {
		case LiveChatMetricKeyAgentsOnline:
			w.AgentsOnline = data
		case LiveChatMetricKeyAgentsAway:
			w.AgentsAway = data
		case LiveChatMetricKeyAgentsInvisible:
			w.AgentsInvisible = data
		}
	}

	w.LastUpdateReceived = &updated
}

func (w *WebsocketAgentMetricData) patchSubscription(update map[LiveChatMetricKeyAgent]bool) {
	for dataPointKey, data := range update {
		switch dataPointKey {
		case LiveChatMetricKeyAgentsOnline:
			w.Subscriptions.AgentsOnline = data
		case LiveChatMetricKeyAgentsAway:
			w.Subscriptions.AgentsAway = data
		case LiveChatMetricKeyAgentsInvisible:
			w.Subscriptions.AgentsInvisible = data
		}
	}
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

	return nil
}

func (s *RealTimeChatStreamingService) ConnectToWebsocket(ctx context.Context) error {
	if err := s.initiateWebsocketConnection(ctx); err != nil {
		return err
	}
	defer s.wsClient.conn.Close()

	reader := concur.NewAsyncReader(
		func(ctx context.Context) (realTimeChatStreamingFrame, error) {
			return s.read()
		},
	)

	go reader.Loop(ctx)
	defer reader.Close()

	keepalive := time.NewTicker(time.Second * 15)
	pongMonitor := time.NewTicker(time.Second * 30)

	for {
		select {
		case update := <-reader.Updates():
			if update.Err != nil {
				return update.Err
			}

			if err := s.handleFrame(update.Item); err != nil {
				return err
			}

		case <-keepalive.C:
			if err := s.ping(); err != nil {
				return err
			}

		case t := <-pongMonitor.C:
			if s.wsCache.metadata.receivedControl == nil {
				if s.wsCache.metadata.connStarted == nil {
					continue
				}

				if t.Sub(*s.wsCache.metadata.connStarted) > (time.Second * 30) {
					return errors.New("could not detect connection started within 30 seconds, killing connection")
				}
			}

			if t.Sub(*s.wsCache.metadata.receivedControl) > (time.Minute * 2) {
				return errors.New("have not received ping within timeframe, killing connection")
			}

		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}

type realTimeChatStreamingFrame struct {
	opCode  ws.OpCode
	payload []byte
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

func (s *RealTimeChatStreamingService) ping() error {
	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpPing,
		payload: nil,
	}); err != nil {
		return err
	}

	t := time.Now()
	s.wsCache.metadata.sentPing = &t

	return nil
}

func (s *RealTimeChatStreamingService) read() (realTimeChatStreamingFrame, error) {
	header, err := ws.ReadHeader(s.wsClient.conn)
	if err != nil {
		return realTimeChatStreamingFrame{}, err
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(s.wsClient.conn, payload)
	if err != nil {
		return realTimeChatStreamingFrame{}, err
	}

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	return realTimeChatStreamingFrame{
		opCode:  header.OpCode,
		payload: payload,
	}, nil
}

func (s *RealTimeChatStreamingService) write(frame realTimeChatStreamingFrame) error {
	if err := s.connectionEstablished(); err != nil {
		return err
	}

	writer := wsutil.NewWriter(s.wsClient.conn, ws.StateClientSide, frame.opCode)
	_, err := writer.Write(frame.payload)
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
	interval := time.NewTicker(time.Millisecond * 250)
	for {
		select {
		case <-timeout.C:
			return fmt.Errorf("timeout reached")
		case <-interval.C:
			if s.wsClient.conn != nil {
				return nil
			}
		}
	}
}

func (s *RealTimeChatStreamingService) handleControlFrame(
	frame realTimeChatStreamingFrame,
) error {
	receivedTime := time.Now()
	s.wsCache.metadata.receivedControl = &receivedTime

	switch frame.opCode {
	case ws.OpPing:
		return s.write(realTimeChatStreamingFrame{
			opCode:  ws.OpPong,
			payload: frame.payload,
		})

	case ws.OpPong:
		return nil

	case ws.OpClose:
		return s.write(realTimeChatStreamingFrame{
			opCode:  ws.OpClose,
			payload: frame.payload,
		})
	}

	return nil
}

func (s *RealTimeChatStreamingService) handleDataFrame(
	frame realTimeChatStreamingFrame,
) error {
	receivedTime := time.Now()
	s.wsCache.metadata.receivedData = &receivedTime

	type status struct {
		StatusCode int `json:"status_code"`
	}

	tempFrame := status{}
	if err := json.Unmarshal(frame.payload, &tempFrame); err != nil {
		return err
	}

	switch tempFrame.StatusCode {
	case http.StatusUnauthorized:
		return ErrRealTimeChatWebsocketUnauthenticated
	case http.StatusInternalServerError:
		return ErrRealTimeChatWebsocketUnsupportedIncomingMessage
	case http.StatusOK:
		// Determine what kind of data frame was sent
		t := map[string]any{}
		if err := json.Unmarshal(frame.payload, &t); err != nil {
			return err
		}

		if content, ok := t["content"]; ok {
			contentBytes, err := json.Marshal(content)
			if err != nil {
				return err
			}

			c := realTimeChatStreamingContent{}
			if err := json.Unmarshal(contentBytes, &c); err != nil {
				return err
			}

			if c.Topic == "agents" {
				data := realTimeChatStreamingContentAgentMetric{}

				if err := json.Unmarshal(contentBytes, &data); err != nil {
					return err
				}

				groupID := realTimeChatStreamingGlobalDepartmentID
				if data.DepartmentID != nil {
					groupID = *data.DepartmentID
				}

				// NOTE: As we are patching the data, it is acceptable to throw away the "ok" value here.
				// The effect is the same - we update a single item on the object and insert it to the cache.
				existingItem, _ := s.wsCache.agent.Get(groupID)

				existingItem.patchData(data.Data)

				s.wsCache.agent.Update(groupID, existingItem)
			}

			if c.Topic == "chats" {
				data := realTimeChatStreamingContentChatMetric{}

				if err := json.Unmarshal(contentBytes, &data); err != nil {
					return err
				}

				groupID := realTimeChatStreamingGlobalDepartmentID
				if data.DepartmentID != nil {
					groupID = *data.DepartmentID
				}

				// NOTE: As we are patching the data, it is acceptable to throw away the "ok" value here.
				// The effect is the same - we update a single item on the object and insert it to the cache.
				existingItem, _ := s.wsCache.chat.Get(groupID)

				existingItem.patchData(data.Data)

				s.wsCache.chat.Update(groupID, existingItem)
			}

			return nil
		}
	}

	return nil
}

type realTimeChatStreamingAgentSubscription struct {
	Topic        string   `json:"topic"`
	Action       string   `json:"action"`
	DepartmentID *GroupID `json:"department_id,omitempty"`
}

type realTimeChatStreamingChatSubscription struct {
	Topic        string              `json:"topic"`
	Action       string              `json:"action"`
	Window       *LiveChatTimeWindow `json:"window,omitempty"`
	DepartmentID *GroupID            `json:"department_id,omitempty"`
}

func (s *RealTimeChatStreamingService) SubscribeToAllAgentMetrics(
	departmentID *GroupID,
) error {
	allMetrics := []LiveChatMetricKeyAgent{
		LiveChatMetricKeyAgentsOnline,
		LiveChatMetricKeyAgentsAway,
		LiveChatMetricKeyAgentsInvisible,
	}

	for _, metric := range allMetrics {
		if err := s.SubscribeToSingleAgentMetric(metric, departmentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *RealTimeChatStreamingService) UnsubscribeFromAllAgentMetrics(
	departmentID *GroupID,
) error {
	allMetrics := []LiveChatMetricKeyAgent{
		LiveChatMetricKeyAgentsOnline,
		LiveChatMetricKeyAgentsAway,
		LiveChatMetricKeyAgentsInvisible,
	}

	for _, metric := range allMetrics {
		if err := s.UnsubscribeFromSingleAgentMetric(metric, departmentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *RealTimeChatStreamingService) SubscribeToAllChatMetrics(
	departmentID *GroupID,
) error {
	allMetrics := []LiveChatMetricKeyChat{
		LiveChatMetricKeyIncomingChats,
		LiveChatMetricKeyAssignedChats,
		LiveChatMetricKeyActiveChats,
		LiveChatMetricKeyWaitingTimeAvg,
		LiveChatMetricKeyWaitingTimeMax,
		LiveChatMetricKeyChatDurationAvg,
		LiveChatMetricKeyChatDurationMax,
		LiveChatMetricKeyResponseTimeAvg,
		LiveChatMetricKeyResponseTimeMax,
	}

	allWindowMetrics := []LiveChatMetricKeyChatWindow{
		LiveChatMetricKeyMissedChats,
		LiveChatMetricKeySatisfactionGood,
		LiveChatMetricKeySatisfactionBad,
	}

	for _, metric := range allMetrics {
		if err := s.SubscribeToSingleChatMetric(metric, departmentID); err != nil {
			return err
		}
	}

	for _, metric := range allWindowMetrics {
		thirty := LiveChatTimeWindow30Minutes
		sixty := LiveChatTimeWindow60Minutes
		if err := s.SubscribeToSingleChatWindowMetric(metric, &thirty, departmentID); err != nil {
			return err
		}

		if err := s.SubscribeToSingleChatWindowMetric(metric, &sixty, departmentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *RealTimeChatStreamingService) UnsubscribeFromAllChatMetrics(
	departmentID *GroupID,
) error {
	if err := s.UnsubscribeFromSingleAgentMetric(LiveChatMetricKeyAgentsOnline, departmentID); err != nil {
		return err
	}

	if err := s.UnsubscribeFromSingleAgentMetric(LiveChatMetricKeyAgentsAway, departmentID); err != nil {
		return err
	}

	return s.UnsubscribeFromSingleAgentMetric(LiveChatMetricKeyAgentsInvisible, departmentID)
}

func (s *RealTimeChatStreamingService) SubscribeToSingleAgentMetric(
	metric LiveChatMetricKeyAgent,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingAgentSubscription{
		Topic:        fmt.Sprintf("agents.%s", metric),
		Action:       "subscribe",
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	existingItem, _ := s.wsCache.agent.Get(groupID)
	existingItem.patchSubscription(map[LiveChatMetricKeyAgent]bool{
		metric: true,
	})

	s.wsCache.agent.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) UnsubscribeFromSingleAgentMetric(
	metric LiveChatMetricKeyAgent,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingAgentSubscription{
		Topic:        fmt.Sprintf("agents.%s", metric),
		Action:       "unsubscribe",
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	existingItem, _ := s.wsCache.agent.Get(groupID)
	existingItem.patchSubscription(map[LiveChatMetricKeyAgent]bool{
		metric: false,
	})

	s.wsCache.agent.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) SubscribeToSingleChatWindowMetric(
	metric LiveChatMetricKeyChatWindow,
	window *LiveChatTimeWindow,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingChatSubscription{
		Topic:        fmt.Sprintf("chats.%s", metric),
		Action:       "subscribe",
		Window:       window,
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	subscriptionKey := string(metric)
	if window != nil {
		subscriptionKey = fmt.Sprintf("%s%d", subscriptionKey, *window)
	}

	existingItem, _ := s.wsCache.chat.Get(groupID)
	existingItem.patchSubscription(map[string]bool{
		subscriptionKey: true,
	})

	s.wsCache.chat.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) SubscribeToSingleChatMetric(
	metric LiveChatMetricKeyChat,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingChatSubscription{
		Topic:        fmt.Sprintf("chats.%s", metric),
		Action:       "subscribe",
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	existingItem, _ := s.wsCache.chat.Get(groupID)
	existingItem.patchSubscription(map[string]bool{
		string(metric): true,
	})

	s.wsCache.chat.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) UnsubscribeFromChatWindowMetric(
	metric LiveChatMetricKeyChatWindow,
	window *LiveChatTimeWindow,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingChatSubscription{
		Topic:        fmt.Sprintf("chats.%s", metric),
		Action:       "unsubscribe",
		Window:       window,
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	subscriptionKey := string(metric)
	if window != nil {
		subscriptionKey = fmt.Sprintf("%s%d", subscriptionKey, *window)
	}

	existingItem, _ := s.wsCache.chat.Get(groupID)
	existingItem.patchSubscription(map[string]bool{
		subscriptionKey: false,
	})

	s.wsCache.chat.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) UnsubscribeFromChatMetric(
	metric LiveChatMetricKeyChat,
	departmentID *GroupID,
) error {
	payload := realTimeChatStreamingChatSubscription{
		Topic:        fmt.Sprintf("chats.%s", metric),
		Action:       "unsubscribe",
		DepartmentID: departmentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.write(realTimeChatStreamingFrame{
		opCode:  ws.OpText,
		payload: payloadBytes,
	}); err != nil {
		return err
	}

	sentDataTime := time.Now()
	s.wsCache.metadata.sentData = &sentDataTime

	groupID := realTimeChatStreamingGlobalDepartmentID
	if departmentID != nil {
		groupID = *departmentID
	}

	existingItem, _ := s.wsCache.chat.Get(groupID)
	existingItem.patchSubscription(map[string]bool{
		string(metric): true,
	})

	s.wsCache.chat.Update(groupID, existingItem)

	return nil
}

func (s *RealTimeChatStreamingService) GetAllAgentMetricsForDepartments() (map[GroupID]WebsocketAgentMetricData, error) {
	items, err := s.wsCache.agent.GetAll()
	if err != nil {
		return nil, err
	}

	return items.Items, nil
}

func (s *RealTimeChatStreamingService) GetAllChatMetricsForDepartments() (map[GroupID]WebsocketChatMetricData, error) {
	items, err := s.wsCache.chat.GetAll()
	if err != nil {
		return nil, err
	}

	return items.Items, nil
}

func (s *RealTimeChatStreamingService) GetAllAgentMetricsByDepartmentID(departmentID GroupID) (WebsocketAgentMetricData, error) {
	item, ok := s.wsCache.agent.Get(departmentID)
	if !ok {
		return WebsocketAgentMetricData{}, errors.New("could not find data for Department")
	}

	return item, nil
}

func (s *RealTimeChatStreamingService) GetAllChatMetricsByDepartmentID(departmentID GroupID) (WebsocketAgentMetricData, error) {
	item, ok := s.wsCache.agent.Get(departmentID)
	if !ok {
		return WebsocketAgentMetricData{}, errors.New("could not find data for Department")
	}

	return item, nil
}

func (s *RealTimeChatStreamingService) GetConnectionStartedTime() *time.Time {
	return s.wsCache.metadata.connStarted
}

type realTimeChatStreamingContent struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

type realTimeChatStreamingContentAgentMetric struct {
	Topic        string                            `json:"topic"`
	Data         map[LiveChatMetricKeyAgent]uint64 `json:"data"`
	Type         string                            `json:"type"`
	DepartmentID *GroupID                          `json:"department_id"`
}

type realTimeChatStreamingContentChatMetric struct {
	Topic        string                        `json:"topic"`
	Data         map[LiveChatMetricKeyChat]any `json:"data"`
	Type         string                        `json:"type"`
	DepartmentID *GroupID                      `json:"department_id"`
}
