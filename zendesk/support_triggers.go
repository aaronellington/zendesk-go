package zendesk

import (
	"context"
	"time"
)

type TriggerResponse struct {
	Trigger Trigger `json:"trigger"`
}

type TriggersResponse struct {
	Triggers []Trigger `json:"triggers"`
	cursorPaginationResponse
}

type Trigger struct {
	ID          TriggerID              `json:"id"`
	URL         string                 `json:"url"`
	Title       string                 `json:"title"`
	Active      bool                   `json:"active"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CreatedAt   time.Time              `json:"created_at"`
	Default     bool                   `json:"default"`
	Actions     []BusinessRuleAction   `json:"actions"`
	Conditions  BusinessRuleConditions `json:"conditions"`
	Description *string                `json:"description"`
	Position    int                    `json:"position"`
	RawTitle    string                 `json:"raw_title"`
	CategoryID  string                 `json:"category_id"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/
type TriggerService struct {
	client  *client
	generic genericService[
		TriggerID,
		TriggerResponse,
		TriggersResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/#list-triggers
func (s *TriggerService) List(
	ctx context.Context,
	pageHandler func(response TriggersResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
