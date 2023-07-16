package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type ChatService struct {
	chatsService  *ChatsService
	agentsService *AgentEventService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *ChatService) Chats() *ChatsService {
	return s.chatsService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/agents/
func (s *ChatService) AgentsService() *AgentEventService {
	return s.agentsService
}
