package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type ChatService struct {
	chatsService       *ChatsService
	agentEventService  *AgentEventService
	departmentService  *DepartmentService
	oauthClientService *OAuthClientService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *ChatService) Chats() *ChatsService {
	return s.chatsService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
func (s *ChatService) AgentEvent() *AgentEventService {
	return s.agentEventService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/departments
func (s *ChatService) Department() *DepartmentService {
	return s.departmentService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/oauth_clients/
func (s *ChatService) OAuthClientService() *OAuthClientService {
	return s.oauthClientService
}
