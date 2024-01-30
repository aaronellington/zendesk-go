package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
type RESTService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#example-response
type ChatMetrics struct {
	MissedChats      *ChatMetricWindow `json:"missed_chats"`
	ChatDurationMax  *uint64           `json:"chat_duration_max"`
	SatisfactionBad  *ChatMetricWindow `json:"satisfaction_bad"`
	ActiveChats      uint64            `json:"active_chats"`
	SatisfactionGood *ChatMetricWindow `json:"satisfaction_good"`
	IncomingChats    uint64            `json:"incoming_chats"`
	AssignedChats    uint64            `json:"assigned_chats"`
	ChatDurationAvg  *uint64           `json:"chat_duration_avg"`
	WaitingTimeAvg   *uint64           `json:"waiting_time_avg"`
	ResponseTimeAvg  *uint64           `json:"response_time_avg"`
	WaitingTimeMax   *uint64           `json:"waiting_time_max"`
	ResponseTimeMax  *uint64           `json:"response_time_max"`
}

type ChatMetricWindow struct {
	SixtyMinuteWindow  uint64 `json:"60"`
	ThirtyMinuteWindow uint64 `json:"30"`
}

type ChatsStreamResponse struct {
	Content    ChatsStreamResponseContent `json:"content"`
	Message    string                     `json:"message"` // https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#data-initialization
	StatusCode int                        `json:"status_code"`
}

type ChatsStreamResponseContent struct {
	Topic        string            `json:"topic"`
	Data         ChatMetrics       `json:"data"`
	Type         string            `json:"type"`
	DepartmentID *ChatDepartmentID `json:"department_id"`
}

func (s *RESTService) GetAllChatMetrics(ctx context.Context) (ChatsStreamResponse, error) {
	return s.getChatMetricsBy(ctx, nil)
}

func (s *RESTService) GetAllChatMetricsByDepartment(ctx context.Context, departmentID ChatDepartmentID) (ChatsStreamResponse, error) {
	filter := url.Values{
		"department": []string{strconv.FormatUint(uint64(departmentID), 10)},
	}
	return s.getChatMetricsBy(ctx, &filter)
}

func (s *RESTService) GetAllChatMetricsForSpecificTimeWindow(ctx context.Context, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error) {
	filter := url.Values{
		"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
	}
	return s.getChatMetricsBy(ctx, &filter)
}

func (s *RESTService) getChatMetricsBy(ctx context.Context, filter *url.Values) (ChatsStreamResponse, error) {
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

func (s *RESTService) GetSingleChatMetric(ctx context.Context, chatMetric LiveChatMetricKeyChat) (ChatsStreamResponse, error) {
	return s.getChatMetricBy(ctx, chatMetric, nil)
}

func (s *RESTService) GetSingleChatMetricForDepartment(ctx context.Context, chatMetric LiveChatMetricKeyChat, departmentID ChatDepartmentID) (ChatsStreamResponse, error) {
	filter := url.Values{
		"department": []string{strconv.FormatUint(uint64(departmentID), 10)},
	}
	return s.getChatMetricBy(ctx, chatMetric, &filter)
}

func (s *RESTService) GetSingleChatMetricForSpecificTimeWindow(ctx context.Context, chatMetric LiveChatMetricKeyChat, timeWindow LiveChatTimeWindow) (ChatsStreamResponse, error) {
	filter := url.Values{
		"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
	}
	return s.getChatMetricBy(ctx, chatMetric, &filter)
}

func (s *RESTService) getChatMetricBy(ctx context.Context, chatMetric LiveChatMetricKeyChat, filter *url.Values) (ChatsStreamResponse, error) {
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
