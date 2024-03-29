package zendesk

type (
	ChatID           uint64
	ChatEngagementID uint64
)

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
type LiveChatChatsService struct {
	c *client
}
