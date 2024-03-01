package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
type RealTimeChatRestService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#example-response
type ChatMetrics struct {
	MissedChats      *ChatMetricWindow `json:"missed_chats"`
	ChatDurationMax  *int64            `json:"chat_duration_max"`
	SatisfactionBad  *ChatMetricWindow `json:"satisfaction_bad"`
	ActiveChats      int64             `json:"active_chats"`
	SatisfactionGood *ChatMetricWindow `json:"satisfaction_good"`
	IncomingChats    int64             `json:"incoming_chats"`
	AssignedChats    int64             `json:"assigned_chats"`
	ChatDurationAvg  *int64            `json:"chat_duration_avg"`
	WaitingTimeAvg   *int64            `json:"waiting_time_avg"`
	ResponseTimeAvg  *int64            `json:"response_time_avg"`
	WaitingTimeMax   *int64            `json:"waiting_time_max"`
	ResponseTimeMax  *int64            `json:"response_time_max"`
}

type ChatMetricWindow struct {
	SixtyMinuteWindow  int64 `json:"60"`
	ThirtyMinuteWindow int64 `json:"30"`
}

type ChatsStreamResponse struct {
	Content    ChatsStreamResponseContent `json:"content"`
	Message    string                     `json:"message"` // https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#data-initialization
	StatusCode int                        `json:"status_code"`
}

type ChatsStreamResponseContent struct {
	Topic        string      `json:"topic"`
	Data         ChatMetrics `json:"data"`
	Type         string      `json:"type"`
	DepartmentID *GroupID    `json:"department_id"`
}

func (s *RealTimeChatRestService) GetAllChatMetrics(ctx context.Context) (ChatsStreamResponse, error) {
	return s.getChatMetricsBy(ctx, nil)
}

func (s *RealTimeChatRestService) GetAllChatMetricsForDepartment(ctx context.Context, departmentID GroupID) (ChatsStreamResponse, error) {
	filter := url.Values{
		"department_id": []string{strconv.FormatUint(uint64(departmentID), 10)},
	}

	return s.getChatMetricsBy(ctx, &filter)
}

func (s *RealTimeChatRestService) GetAllChatMetricsForSpecificTimeWindow(ctx context.Context, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error) {
	filter := url.Values{
		"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
	}

	return s.getChatMetricsBy(ctx, &filter)
}

func (s *RealTimeChatRestService) getChatMetricsBy(ctx context.Context, filter *url.Values) (ChatsStreamResponse, error) {
	filters := ""
	if filter != nil {
		filters = fmt.Sprintf("?%s", filter.Encode())
	}

	target := ChatsStreamResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/stream/chats%s", filters),
		http.NoBody,
	)
	if err != nil {
		return ChatsStreamResponse{}, err
	}

	if err := s.client.RealTimeChatRequest(request, &target); err != nil {
		return ChatsStreamResponse{}, err
	}

	return target, nil
}

func (s *RealTimeChatRestService) GetSingleChatMetric(ctx context.Context, chatMetric LiveChatMetricKeyChat) (ChatsStreamResponse, error) {
	return s.getChatMetricBy(ctx, chatMetric, nil)
}

func (s *RealTimeChatRestService) GetSingleChatMetricForDepartment(ctx context.Context, chatMetric LiveChatMetricKeyChat, departmentID GroupID) (ChatsStreamResponse, error) {
	filter := url.Values{
		"department_id": []string{strconv.FormatUint(uint64(departmentID), 10)},
	}

	return s.getChatMetricBy(ctx, chatMetric, &filter)
}

func (s *RealTimeChatRestService) GetSingleChatMetricForSpecificTimeWindow(ctx context.Context, chatMetric LiveChatMetricKeyChat, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error) {
	filter := url.Values{
		"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
	}

	return s.getChatMetricBy(ctx, chatMetric, &filter)
}

func (s *RealTimeChatRestService) getChatMetricBy(ctx context.Context, chatMetric LiveChatMetricKeyChat, filter *url.Values) (ChatsStreamResponse, error) {
	filters := ""
	if filter != nil {
		filters = fmt.Sprintf("?%s", filter.Encode())
	}

	target := ChatsStreamResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/stream/chats/%s%s", chatMetric, filters),
		http.NoBody,
	)
	if err != nil {
		return ChatsStreamResponse{}, err
	}

	if err := s.client.RealTimeChatRequest(request, &target); err != nil {
		return ChatsStreamResponse{}, err
	}

	return target, nil
}
