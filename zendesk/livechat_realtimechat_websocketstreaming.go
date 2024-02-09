package zendesk

import (
	"context"
	"fmt"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
type RealTimeChatStreamingService struct {
	client *client
	wsConn *net.Conn
}

func (s *RealTimeChatStreamingService) InitiateWebsocketConnection(ctx context.Context) error {
	// Connection has been made already
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

	conn, buff, hs, err := dialer.Dial(ctx, "wss://rtm.zopim.com/stream")
	if err != nil {
		return err
	}

	for conn != nil {
		msg, _, err := wsutil.ReadServerData(conn)
		if err != nil {
			return err
		}
		fmt.Println(string(msg))
	}

	fmt.Println(buff, hs)

	s.wsConn = &conn

	return nil
}
