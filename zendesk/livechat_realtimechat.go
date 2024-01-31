package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type RealTimeChatService struct {
	realTimeChatRestService      *RealTimeChatRestService
	realTimeChatStreamingService *RealTimeChatStreamingService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
func (s *RealTimeChatService) REST() *RealTimeChatRestService {
	return s.realTimeChatRestService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
func (s *RealTimeChatService) WebsocketStreaming() *RealTimeChatStreamingService {
	return s.realTimeChatStreamingService
}
