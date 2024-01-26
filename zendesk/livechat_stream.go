package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

type ChatStreamService struct {
	client *client
}

type ChatsStreamResponse struct {
	ChatDurationMax *int `json:"chat_duration_max"`
	ActiveChats     *int `json:"active_chats"`
	IncomingChats   *int `json:"incoming_chats"`
	AssignedChats   *int `json:"assigned_chats"`
	ChatDurationAvg *int `json:"chat_duration_avg"`
	WaitingTimeAvg  *int `json:"waiting_time_avg"`
	ResponseTimeAvg *int `json:"response_time_avg"`
	WaitingTimeMax  *int `json:"waiting_time_max"`
	ResponseTimeMax *int `json:"response_time_max"`
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#data-initialization

func (s *ChatStreamService) List(ctx context.Context, departmentID string, pageHandler func(page ChatsStreamResponse) error) error {
	// baseURL := "https://rtm.zopim.com/stream/chats?department_id=%s"
	requestURL := fmt.Sprintf("https://rtm.zopim.com/stream/chats?department_id=%s", departmentID)

	for {
		target := ChatsStreamResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			requestURL,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.LiveChatRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

	}

}
