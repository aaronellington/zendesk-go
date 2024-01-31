package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type RealTimeChatService struct {
	realTimeChatRestService      *RealTimeChatRestService
	realTimeChatStreamingService *RealTimeChatStreamingService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
func (s *RealTimeChatService) RealTimeChatRestService() *RealTimeChatRestService {
	return s.realTimeChatRestService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
func (s *RealTimeChatService) RealTimeChatStreamingService() *RealTimeChatStreamingService {
	return s.realTimeChatStreamingService
}
