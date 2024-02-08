package zendesk

// https://developer.zendesk.com/api-reference/live-chat/introduction/
type LiveChatService struct {
	chatService              *ChatService
	realTimeChatService      *RealTimeChatService
	chatConversationsService *ChatConversationsService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
func (s *LiveChatService) Chat() *ChatService {
	return s.chatService
}

// https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/introduction/
func (s *LiveChatService) RealTimeChat() *RealTimeChatService {
	return s.realTimeChatService
}

// https://developer.zendesk.com/api-reference/live-chat/chat-conversations-api/conversations-api/
func (s *LiveChatService) ChatConversations() *ChatConversationsService {
	return s.chatConversationsService
}
