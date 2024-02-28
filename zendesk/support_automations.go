package zendesk

import (
	"context"
	"time"
)

type AutomationResponse struct {
	Automation Automation `json:"automation"`
}

type AutomationsResponse struct {
	Automations []Automation `json:"automations"`
	cursorPaginationResponse
}

type Automation struct {
	ID         AutomationID           `json:"id"`
	URL        string                 `json:"url"`
	Title      string                 `json:"title"`
	Active     bool                   `json:"active"`
	UpdatedAt  time.Time              `json:"updated_at"`
	CreatedAt  time.Time              `json:"created_at"`
	Default    bool                   `json:"default"`
	Actions    []BusinessRuleAction   `json:"actions"`
	Conditions BusinessRuleConditions `json:"conditions"`
	Position   int                    `json:"position"`
	RawTitle   string                 `json:"raw_title"`
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/
type AutomationService struct {
	client  *client
	generic genericService[
		AutomationID,
		AutomationResponse,
		AutomationsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/#list-automations
func (s AutomationService) List(
	ctx context.Context,
	pageHandler func(response AutomationsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
