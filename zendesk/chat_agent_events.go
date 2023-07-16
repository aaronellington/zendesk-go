package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type AgentEventExportResponse struct {
	AgentEvents []AgentEvent `json:"agent_events"`
	EndTimeUnix float64      `json:"end_time"`
	NextPage    string       `json:"next_page"`
	Count       int64        `json:"count"`
}

func (response AgentEventExportResponse) EndTime() time.Time {
	return time.Unix(int64(response.EndTimeUnix), 0)
}

type AgentEvent struct {
	StartTime     time.Time       `json:"timestamp"`
	AccountID     uint64          `json:"account_id"`
	AgentID       UserID          `json:"agent_id"`
	FieldName     string          `json:"field_name"`
	ID            string          `json:"id"`
	PreviousValue AgentEventValue `json:"previous_value"`
	Value         AgentEventValue `json:"value"`
}

type AgentEventValue string

func (a *AgentEventValue) UnmarshalJSON(data []byte) (err error) {
	var intTarget int64
	if err := json.Unmarshal(data, &intTarget); err == nil {
		*a = AgentEventValue(fmt.Sprintf("%d", intTarget))
		return nil
	}

	var stringTarget string

	if err := json.Unmarshal(data, &stringTarget); err == nil {
		*a = AgentEventValue(stringTarget)
		return nil
	}

	return errors.New("invalid type")
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/
type AgentEventService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/live-chat/chat-api/incremental_agent_events_api/#incremental-agent-events-export
func (s *AgentEventService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(response AgentEventExportResponse) error,
) error {
	const limit = 1000
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime.Unix()))
	query.Set("limit", fmt.Sprintf("%d", limit))
	url := fmt.Sprintf("/api/v2/incremental/agent_events?%s", query.Encode())

	for {
		target := AgentEventExportResponse{}

		if err := s.client.ChatRequest(
			ctx,
			http.MethodGet,
			url,
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

		url = target.NextPage
	}

	return nil
}
