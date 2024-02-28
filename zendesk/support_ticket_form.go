package zendesk

import (
	"context"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
type TicketFormService struct {
	client  *client
	generic genericService[
		TicketFormID,
		TicketFormResponse,
		TicketFormsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#json-format
type TicketForm struct {
	URL            string          `json:"url"`
	Name           string          `json:"name"`
	DisplayName    string          `json:"display_name"`
	ID             TicketFormID    `json:"id"`
	RawName        string          `json:"raw_name"`
	RawDisplayName string          `json:"raw_display_name"`
	EndUserVisible bool            `json:"end_user_visible"`
	Position       int             `json:"position"`
	TicketFieldIds []TicketFieldID `json:"ticket_field_ids"`
	Active         bool            `json:"active"`
	Default        bool            `json:"default"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	InAllBrands    bool            `json:"in_all_brands"`
}

type TicketFormsResponse struct {
	TicketForms []TicketForm `json:"ticket_forms"`
	offsetPaginationResponse
}

type TicketFormResponse struct {
	TicketForm TicketForm `json:"ticket_form"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#list-ticket-forms
func (s TicketFormService) List(
	ctx context.Context,
	pageHandler func(response TicketFormsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#show-ticket-form
func (s TicketFormService) Show(
	ctx context.Context,
	id TicketFormID,
) (TicketFormResponse, error) {
	return s.generic.Show(ctx, id)
}
