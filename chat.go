package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type ChatService struct {
	chatsService *ChatsService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *ChatService) Chats() *ChatsService {
	return s.chatsService
}
