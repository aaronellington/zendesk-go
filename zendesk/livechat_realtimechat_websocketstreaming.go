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
	individualDepartments *utils.MemoryCacheInstance[GroupID, WebsocketAgentMetricData]
	globalData            WebsocketAgentMetricData
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
	MissedChats        *ChatMetricWindow              `json:"missed_chats"`
	ChatDurationMax    *uint64                        `json:"chat_duration_max"`
	SatisfactionBad    *ChatMetricWindow              `json:"satisfaction_bad"`
	ActiveChats        uint64                         `json:"active_chats"`
	SatisfactionGood   *ChatMetricWindow              `json:"satisfaction_good"`
	IncomingChats      uint64                         `json:"incoming_chats"`
	AssignedChats      uint64                         `json:"assigned_chats"`
	ChatDurationAvg    *uint64                        `json:"chat_duration_avg"`
	WaitingTimeAvg     *uint64                        `json:"waiting_time_avg"`
	ResponseTimeAvg    *uint64                        `json:"response_time_avg"`
	WaitingTimeMax     *uint64                        `json:"waiting_time_max"`
	ResponseTimeMax    *uint64                        `json:"response_time_max"`
	Subscriptions      map[LiveChatMetricKeyChat]bool `json:"subscriptions"`
	LastUpdateReceived *time.Time                     `json:"last_update_received"`
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
	Data               map[LiveChatMetricKeyAgent]uint64 `json:"data"`
	Subscriptions      map[LiveChatMetricKeyAgent]bool   `json:"subscriptions"`
	LastUpdateReceived *time.Time                        `json:"last_update_received"`
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

	keepalive := time.NewTicker(time.Second * 1)
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

func (s *RealTimeChatStreamingService) CloseConnection(closeStatusCode WebsocketCloseCode, closeReason string) error {
	closeFrame := RealTimeChatStreamingCloseFrame{
		Code:   closeStatusCode,
		Reason: closeReason,
	}

	payload, err := json.Marshal(closeFrame)
	if err != nil {
		return err
	}

	if err := s.write(
		realTimeChatStreamingFrame{
			opCode:  ws.OpClose,
			payload: payload,
		},
	); err != nil {
		return err
	}

	return nil
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
		return errors.New("error, unsupported message sent to server: " + string(frame.payload))
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

			c := RealTimeChatStreamingContent{}
			if err := json.Unmarshal(contentBytes, &c); err != nil {
				return err
			}

			if c.Topic == "agents" {
				data := RealTimeChatStreamingContentAgentMetric{}

				if err := json.Unmarshal(contentBytes, &data); err != nil {
					return err
				}

				if data.DepartmentID != nil {
					oldData, ok := s.wsCache.agent.individualDepartments.Get(*data.DepartmentID)
					if !ok {
						// s.wsCache.agent.individualDepartments.Update(
						// 	*data.DepartmentID, WebsocketAgentMetricData{
						// 		Subscriptions: make(map[LiveChatMetricKeyAgent]bool),
						// 		Data:          data.Data,
						// 	},
						// )

						return errors.New("Ahhh!")
					}

					for k, v := range data.Data {
						oldData.Data[k] = v
					}

					s.wsCache.agent.individualDepartments.Update(*data.DepartmentID, oldData)

				}
			}

			return nil
		}

	}

	return nil
}

type GlobalSubscription struct {
	Topic  string `json:"topic"`
	Action string `json:"action"`
}

type DepartmentSubscription struct {
	Topic        string  `json:"topic"`
	Action       string  `json:"action"`
	DepartmentID GroupID `json:"department_id"`
}

func (s *RealTimeChatStreamingService) SubscribeToAgentMetricByDepartment(ctx context.Context, metric LiveChatMetricKeyAgent, departmentID GroupID) error {
	payload := DepartmentSubscription{
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

	oldData, ok := s.wsCache.agent.individualDepartments.Get(departmentID)
	if !ok {
		s.wsCache.agent.individualDepartments.Update(
			departmentID, WebsocketAgentMetricData{
				Subscriptions: map[LiveChatMetricKeyAgent]bool{
					metric: true,
				},
				Data: make(map[LiveChatMetricKeyAgent]uint64),
			},
		)

		return nil
	}

	oldData.Subscriptions[metric] = true

	return nil
}

func (s *RealTimeChatStreamingService) SubscribeToAgentMetricGlobal(ctx context.Context, metric LiveChatMetricKeyAgent) error {
	payload := DepartmentSubscription{
		Topic:        fmt.Sprintf("agents.%s", metric),
		Action:       "subscribe",
		DepartmentID: GroupID(13388700431505),
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
	s.wsCache.agent.globalData.Subscriptions[metric] = true

	return nil
}

func (s *RealTimeChatStreamingService) GetAllAgentMetricsForDepartments(ctx context.Context) (map[GroupID]WebsocketAgentMetricData, error) {
	items, err := s.wsCache.agent.individualDepartments.GetAll()
	if err != nil {
		return nil, err
	}

	return items.Items, nil
}

func (s *RealTimeChatStreamingService) GetAllAgentMetricsByDepartmentID(ctx context.Context, departmentID GroupID) (WebsocketAgentMetricData, error) {
	item, ok := s.wsCache.agent.individualDepartments.Get(departmentID)
	if !ok {
		return WebsocketAgentMetricData{}, errors.New("could not find data for Department")
	}

	return item, nil
}

func (s *RealTimeChatStreamingService) GetAllAgentMetricsGlobal(ctx context.Context) (WebsocketAgentMetricData, error) {
	if s.wsCache.agent.globalData.LastUpdateReceived == nil {
		return WebsocketAgentMetricData{}, errors.New("no update received from Zendesk for metric")
	}

	return s.wsCache.agent.globalData, nil
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

type RealTimeChatStreamingContent struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

type RealTimeChatStreamingContentAgentMetric struct {
	Topic        string                            `json:"topic"`
	Data         map[LiveChatMetricKeyAgent]uint64 `json:"data"`
	Type         string                            `json:"type"`
	DepartmentID *GroupID                          `json:"department_id"`
}

type RealTimeChatStreamingContentChatMetric struct {
	Topic        string                           `json:"topic"`
	Data         map[LiveChatMetricKeyChat]uint64 `json:"data"`
	Type         string                           `json:"type"`
	DepartmentID *GroupID                         `json:"department_id"`
}

type RealTimeChatStreamingCloseFrame struct {
	Code   WebsocketCloseCode `json:"code"`
	Reason string             `json:"reason,omitempty"`
}
