package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ChatID string

type ChatsResponse struct {
	Chats   []Chat  `json:"chats"`
	NextURL *string `json:"next_url"`
}

type ChatsIncrementalExportResponse struct {
	Chats    []Chat `json:"chats"`
	Count    uint64 `json:"count"`
	EndID    ChatID `json:"end_id"`
	EndTime  uint64 `json:"end_time"`
	NextPage string `json:"next_page"`
}

type Chat struct {
	ID              ChatID           `json:"id"`
	Visitor         ChatVisitor      `json:"visitor"`
	StartedBy       string           `json:"started_by"`
	Session         ChatSession      `json:"session"`
	WebPath         []ChatWebPath    `json:"webpath"`
	Timestamp       time.Time        `json:"timestamp"`
	Count           ChatCount        `json:"count"`
	Duration        int              `json:"duration"`
	ResponseTime    ChatResponseTime `json:"response_time"`
	AgentIds        []UserID         `json:"agent_ids"`
	Triggered       bool             `json:"triggered"`
	Unread          bool             `json:"unread"`
	Missed          bool             `json:"missed"`
	Tags            []string         `json:"tags"`
	Type            string           `json:"type"`
	History         []ChatHistory    `json:"history"`
	DepartmentID    GroupID          `json:"department_id"`
	EndTimestamp    time.Time        `json:"end_timestamp"`
	ZendeskTicketID TicketID         `json:"zendesk_ticket_id"`
}

type ChatVisitor struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type ChatSession struct {
	EndDate     time.Time `json:"end_date"`
	CountryCode string    `json:"country_code"`
	City        string    `json:"city"`
	Browser     string    `json:"browser"`
	IP          string    `json:"ip"`
	CountryName string    `json:"country_name"`
	ID          string    `json:"id"`
	Region      string    `json:"region"`
	Platform    string    `json:"platform"`
	UserAgent   string    `json:"user_agent"`
	StartDate   time.Time `json:"start_date"`
}

type ChatWebPath struct {
	Timestamp time.Time `json:"timestamp"`
	To        string    `json:"to"`
	From      string    `json:"from"`
	Title     string    `json:"title"`
}

type ChatCount struct {
	Visitor uint `json:"visitor"`
	Agent   uint `json:"agent"`
	Total   uint `json:"total"`
}

type ChatResponseTime struct {
	First uint64  `json:"first"`
	Avg   float64 `json:"avg"`
	Max   uint64  `json:"max"`
}

type ChatHistory struct {
	DepartmentID   GroupID   `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Name           string    `json:"name"`
	Channel        string    `json:"channel"`
	Index          int       `json:"index"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
	Msg            string    `json:"msg"`
	Options        string    `json:"options"`
	MsgID          string    `json:"msg_id"`
	SenderType     string    `json:"sender_type"`
	Source         string    `json:"source"`
	AgentID        UserID    `json:"agent_id,string"`
	Reason         string    `json:"reason"`
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/
type ChatsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/#list-chats
func (s *ChatsService) List(ctx context.Context, pageHandler func(page ChatsResponse) error) error {
	for {
		target := ChatsResponse{}

		if err := s.c.ChatRequest(ctx, http.MethodGet, "/api/v2/chats", http.NoBody, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.NextURL == nil {
			break
		}
	}

	return nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/#show-chat
func (s *ChatsService) Show(ctx context.Context, id ChatID) (Chat, error) {
	target := Chat{}

	if err := s.c.ChatRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/chats/%s", id),
		http.NoBody,
		&target,
	); err != nil {
		return Chat{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_export/#incremental-chat-export
func (s *ChatsService) IncrementalExport(
	ctx context.Context,
	startTime int64,
	pageHandler func(response ChatsIncrementalExportResponse) error,
) error {
	const limit = 1000
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime))
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("fields", "chats(*)")

	for {
		target := ChatsIncrementalExportResponse{}

		if err := s.c.ChatRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/chats?%s", query.Encode()),
			http.NoBody,
			&target,
		); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.Count < limit {
			break
		}

		query.Set("start_time", fmt.Sprintf("%d", target.EndTime))
		query.Set("start_id", string(target.EndID))
	}

	return nil
}
