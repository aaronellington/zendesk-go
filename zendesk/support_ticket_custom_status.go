package zendesk

import (
	"context"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/
type CustomStatusService struct {
	client  *client
	generic genericService[
		CustomStatusID,
		CustomStatusResponse,
		CustomStatusesResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#json-format
type CustomStatus struct {
	Active         bool                 `json:"active"`
	ID             CustomStatusID       `json:"id"`
	AgentLabel     string               `json:"agent_label"`
	StatusCategory CustomStatusCategory `json:"status_category"`
}

type CustomStatusCategory string

const (
	StatusCategoryNew     CustomStatusCategory = "new"
	StatusCategoryOpen    CustomStatusCategory = "open"
	StatusCategoryPending CustomStatusCategory = "pending"
	StatusCategoryHold    CustomStatusCategory = "hold"
	// Tickets with a "Closed" status belong to the "StatusCategorySolved" status category.
	StatusCategorySolved CustomStatusCategory = "solved"
)

type CustomStatusesResponse struct {
	CustomStatuses []CustomStatus `json:"custom_statuses"`
	offsetPaginationResponse
}

type CustomStatusResponse struct {
	CustomStatus CustomStatus `json:"custom_status"`
}

type CustomStatusPayload struct {
	CustomStatus any `json:"custom_status"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#list-custom-ticket-statuses
func (s CustomStatusService) List(
	ctx context.Context,
	pageHandler func(response CustomStatusesResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#show-custom-ticket-status
func (s CustomStatusService) Show(ctx context.Context, id CustomStatusID) (CustomStatusResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#create-custom-ticket-status
func (s CustomStatusService) Create(ctx context.Context, payload CustomStatusPayload) (CustomStatusResponse, error) {
	return s.generic.Create(ctx, payload)
}
