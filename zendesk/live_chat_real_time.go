package zendesk

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

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/introduction/
type LiveChatRealTimeService struct {
	restService      *LiveChatRealTimeRESTService
	streamingService *LiveChatRealTimeStreamingService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
func (s *LiveChatRealTimeService) REST() *LiveChatRealTimeRESTService {
	return s.restService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
func (s *LiveChatRealTimeService) Streaming() *LiveChatRealTimeStreamingService {
	return s.streamingService
}
