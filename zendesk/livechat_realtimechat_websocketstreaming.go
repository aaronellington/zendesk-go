package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
type RealTimeChatStreamingService struct {
	client *client
	wsConn *net.Conn
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
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := encoder.Encode(nil); err != nil {
				return err
			}

			if err := writer.Flush(); err != nil {
				return err
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *RealTimeChatStreamingService) Read(ctx context.Context) error {
	reader := wsutil.NewClientSideReader(*s.wsConn)
	decoder := json.NewDecoder(reader)
	for {
		header, err := reader.NextFrame()
		if err != nil {
			return err
		}

		if header.OpCode == ws.OpClose {
			return io.EOF
		}

		if header.OpCode == ws.OpPong {
			fmt.Println("It")
			if err := reader.Discard(); err != nil {
				return err
			}
		} else {

			b := map[string]any{}
			if err := decoder.Decode(&b); err != nil {
				return err
			}

			fmt.Println(b)
		}
	}
}

func (s *RealTimeChatStreamingService) write(ctx context.Context) error {
	writer := wsutil.NewWriter(*s.wsConn, ws.StateClientSide, ws.OpPing)
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
	// return nil
}

func (s *RealTimeChatRestService) handleControlFrame(ctx context.Context) error {
	return nil
}

type Subscription struct {
	Topic  string `json:"topic"`
	Action string `json:"action"`
}

func (s *RealTimeChatStreamingService) SubscribeToAgentMetric()                     {}
func (s *RealTimeChatStreamingService) SubscribeToChatMetric()                      {}
func (s *RealTimeChatStreamingService) SubscribeToChatMetricForSpecificTimeWindow() {}
