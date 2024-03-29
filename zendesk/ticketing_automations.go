package zendesk

import (
	"context"
	"time"
)

type AutomationID uint64

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/#json-format
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
type TicketingAutomationsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/#show-automation
func (s *TicketingAutomationsService) Show(
	ctx context.Context,
	id AutomationID,
) (AutomationsResponse, error) {
	return showRequest[AutomationID, AutomationsResponse](
		ctx,
		s.c,
		id,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/#list-automations
func (s *TicketingAutomationsService) List(
	ctx context.Context,
	pageHandler func(response AutomationsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(
		ctx,
		s.c,
		pageHandler,
		requestQueryModifiers...,
	)
}

type AutomationResponse struct {
	Automation Automation `json:"automation"`
	ticketingAutomationsObject
}

type AutomationsResponse struct {
	Automations []Automation `json:"automations"`
	ticketingAutomationsObject
	cursorPaginationResponse
}

type ticketingAutomationsObject struct{}

func (r ticketingAutomationsObject) zendeskEntityName() string {
	return "automations"
}
