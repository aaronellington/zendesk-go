package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#get-chat-metrics-by-department
type RealTimeService struct {
	client *client
}

type ChatMetrics struct {
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

type ChatsStreamResponse struct {
	Content    ChatsStreamResponseContent `json:"content"`
	StatusCode int                        `json:"status_code"`
}

type ChatsStreamResponseContent struct {
	Topic        string      `json:"topic"`
	Data         ChatMetrics `json:"data"`
	Type         string      `json:"type"`
	DepartmentID int64       `json:"department_id"`
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#get-chat-metrics-by-department
func (s *RealTimeService) GetChatMetricsByDepartment(ctx context.Context, departmentID int64) (ChatsStreamResponse, error) {
	requestURL := fmt.Sprintf("https://rtm.zopim.com/stream/chats?department_id=%d", departmentID)
	target := ChatsStreamResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
		http.NoBody,
	)
	if err != nil {
		return ChatsStreamResponse{}, err
	}

	if err := s.client.LiveChatRequest(request, &target); err != nil {
		return ChatsStreamResponse{}, err
	}

	return target, nil
}
