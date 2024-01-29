package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type LiveChatService struct {
	chatService       *ChatService
	agentEventService *AgentEventService
	realTimeService   *RealTimeService
	departmentService *DepartmentService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *LiveChatService) Chat() *ChatService {
	return s.chatService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
func (s *LiveChatService) AgentEvent() *AgentEventService {
	return s.agentEventService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/#get-chat-metrics-by-department
func (s *LiveChatService) RealTime() *RealTimeService {
	return s.realTimeService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments
func (s *LiveChatService) Department() *DepartmentService {
	return s.departmentService
}
