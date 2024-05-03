package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ChatsService struct {
	client *client
}

type ChatsResponse struct {
	Chats   []Chat  `json:"chats"`
	NextURL *string `json:"next_url"`
}

type ChatsSearchResponse struct {
	Results []ChatsSearchResult `json:"results"`
	NextURL *string             `json:"next_url"`
}

type ChatsSearchResult struct {
	ID        ChatID    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Preview   string    `json:"preview"`
	Type      string    `json:"type"`
}

type ChatsIncrementalExportResponse struct {
	Chats       []IncrementalExportChat `json:"chats"`
	Count       uint64                  `json:"count"`
	EndID       ChatID                  `json:"end_id"`
	EndTimeUnix int64                   `json:"end_time"`
	NextPage    string                  `json:"next_page"`
}

func (response ChatsIncrementalExportResponse) EndTime() time.Time {
	return time.Unix(response.EndTimeUnix, 0)
}

type IncrementalExportChat struct {
	Chat
	ChatEngagements []ChatEngagement `json:"engagements"`
}

type Chat struct {
	ID        ChatID      `json:"id"`
	Visitor   ChatVisitor `json:"visitor"`
	StartedBy string      `json:"started_by"`
	Session   ChatSession `json:"session"`
	// WebPath         []ChatWebPath    `json:"webpath"`
	Timestamp       time.Time        `json:"timestamp"`
	Count           ChatCount        `json:"count"`
	Duration        uint64           `json:"duration"`
	ResponseTime    ChatResponseTime `json:"response_time"`
	AgentIds        []UserID         `json:"agent_ids"`
	Triggered       bool             `json:"triggered"`
	Unread          bool             `json:"unread"`
	Missed          bool             `json:"missed"`
	Tags            []Tag            `json:"tags"`
	Type            string           `json:"type"`
	History         []ChatHistory    `json:"history"`
	DepartmentID    *GroupID         `json:"department_id"`
	EndTimestamp    time.Time        `json:"end_timestamp"`
	ZendeskTicketID TicketID         `json:"zendesk_ticket_id"`
}

type ChatEngagement struct {
	ID           ChatEngagementID `json:"id"`
	AgentID      UserID           `json:"agent_id"`
	DepartmentID *GroupID         `json:"department_id"`
	Assigned     bool             `json:"assigned"`
	Accepted     bool             `json:"accepted"`
	StartedBy    string           `json:"started_by"`
	Timestamp    time.Time        `json:"timestamp"`
	Duration     float64          `json:"duration"`
	Count        ChatCount        `json:"count"`
	ResponseTime ChatResponseTime `json:"response_time"`
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
	Visitor uint64 `json:"visitor"`
	Agent   uint64 `json:"agent"`
	Total   uint64 `json:"total"`
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

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/#list-chats
func (s *ChatsService) List(ctx context.Context, pageHandler func(page ChatsResponse) error) error {
	requestURL := "/api/v2/chats"

	for {
		target := ChatsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			requestURL,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ChatRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.NextURL == nil {
			break
		}

		requestURL = *target.NextURL
	}

	return nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/#search-chats
func (s *ChatsService) Search(
	ctx context.Context,
	query string,
	pageHandler func(page ChatsSearchResponse) error,
) error {
	values := &url.Values{}
	values.Set("q", query)

	requestURL := "/api/v2/chats/search?" + values.Encode()

	for {
		target := ChatsSearchResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			requestURL,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ChatRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.NextURL == nil {
			break
		}

		requestURL = *target.NextURL
	}

	return nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/chats/#show-chat
func (s *ChatsService) Show(ctx context.Context, id ChatID) (Chat, error) {
	target := Chat{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/chats/%s", id),
		http.NoBody,
	)
	if err != nil {
		return Chat{}, err
	}

	if err := s.client.ChatRequest(request, &target); err != nil {
		return Chat{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_export/#incremental-chat-export
func (s *ChatsService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(response ChatsIncrementalExportResponse) error,
) error {
	const limit = 1000

	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime.Unix()))
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("fields", "chats(*)")

	for {
		target := ChatsIncrementalExportResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/chats?%s", query.Encode()),
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ChatRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.Count < limit {
			break
		}

		query.Set("start_time", fmt.Sprintf("%d", target.EndTimeUnix))
		query.Set("start_id", string(target.EndID))
	}

	return nil
}
