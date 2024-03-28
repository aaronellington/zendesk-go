package zendesk

import "time"

type AgentState struct {
	AgentID         UserID
	EngagementCount uint64
	Status          string
	StatusSince     time.Time
	Timestamp       time.Time
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
type LiveChatIncrementalAgentEventsService struct {
	c *client
}
