package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type LiveChatService struct {
	chatService       *ChatService
	agentEventService *AgentEventService
	chatStreamService *ChatStreamService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *LiveChatService) Chat() *ChatService {
	return s.chatService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
func (s *LiveChatService) AgentEvent() *AgentEventService {
	return s.agentEventService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
func (s *LiveChatService) ChatStream() *ChatStreamService {
	return s.chatStreamService
}
