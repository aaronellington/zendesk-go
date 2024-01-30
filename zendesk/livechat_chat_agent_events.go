package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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
	AccountID     ChatAccountID   `json:"account_id"`
	AgentID       UserID          `json:"agent_id"`
	FieldName     string          `json:"field_name"`
	ID            string          `json:"id"`
	PreviousValue AgentEventValue `json:"previous_value"`
	Value         AgentEventValue `json:"value"`
}

type AgentEventValue string

func (a *AgentEventValue) UnmarshalJSON(data []byte) error {
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
	client               *client
	agentStatesMutex     *sync.Mutex
	agentStates          AgentStates
	agentStatesStartTime *time.Time
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
	url := fmt.Sprintf("https://www.zopim.com/api/v2/incremental/agent_events?%s", query.Encode())

	for {
		target := AgentEventExportResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			url,
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

		url = target.NextPage
	}

	return nil
}

type AgentStates map[UserID]AgentState

type AgentState struct {
	AgentID         UserID
	EngagementCount uint64
	Status          string
	StatusSince     time.Time
	Timestamp       time.Time
}

func (s *AgentEventService) GetAgentStates(
	ctx context.Context,
) AgentStates {
	out := AgentStates{}

	agentStateBytes, err := json.Marshal(s.agentStates)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(agentStateBytes, &out); err != nil {
		panic(err)
	}

	return out
}

func (s *AgentEventService) UpdateAgentStates(
	ctx context.Context,
	defaultStateTime time.Time,
) error {
	if s.agentStatesStartTime == nil {
		s.agentStatesStartTime = &defaultStateTime
	}

	return s.IncrementalExport(
		ctx, *s.agentStatesStartTime,
		func(response AgentEventExportResponse) error {
			s.agentStatesMutex.Lock()
			defer s.agentStatesMutex.Unlock()

			for _, agentEvent := range response.AgentEvents {
				agentState := s.agentStates[agentEvent.AgentID]

				agentState.AgentID = agentEvent.AgentID
				agentState.Timestamp = agentEvent.StartTime

				// Fix an issue where an agent has not changed their status before
				// the defaultStateTime so we don't ever get a status update.
				// But we do get engagements so we should just default the status fields
				// to something reasonable so they are at least set
				{
					if agentState.StatusSince.IsZero() {
						agentState.StatusSince = agentEvent.StartTime
					}

					if agentState.Status == "" {
						agentState.Status = "unknown"
					}
				}

				switch agentEvent.FieldName {
				case "engagements":
					engagementCount, err := strconv.ParseUint(string(agentEvent.Value), 10, 64)
					if err != nil {
						return err
					}

					agentState.EngagementCount = engagementCount
				case "status":
					if agentEvent.Value == "offline" || agentEvent.Value == "invisible" {
						delete(s.agentStates, agentEvent.AgentID)

						continue
					}

					agentState.Status = string(agentEvent.Value)
					agentState.StatusSince = agentEvent.StartTime
				}

				s.agentStates[agentEvent.AgentID] = agentState
			}

			endTime := response.EndTime()
			s.agentStatesStartTime = &endTime

			return nil
		},
	)
}
