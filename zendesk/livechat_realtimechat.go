package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type RealTimeChatService struct {
	restService               *RESTService
	websocketStreamingService *WebsocketStreamingService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
func (s *RealTimeChatService) REST() *RESTService {
	return s.restService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/streaming/
func (s *RealTimeChatService) WebsocketStreaming() *WebsocketStreamingService {
	return s.websocketStreamingService
}
