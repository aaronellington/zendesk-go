package zendesk

import (
	"context"
)

type ticketingTicketFormObject struct{}

func (r ticketingTicketFormObject) zendeskEntityName() string {
	return "ticket_forms"
}

type TicketFormID int64

type TicketForm struct {
	ID TicketFormID `json:"id"`
}

type TicketFormResponse struct {
	TicketForm TicketForm `json:"ticket_form"`
	ticketingTicketFormObject
}

type TicketFormsResponse struct {
	TicketForms []TicketForm `json:"ticket_forms"`
	ticketingTicketFormObject
	cursorPaginationResponse
}

type TicketFormPayload struct {
	TicketForm any `json:"ticket_form"`
}

type TicketFormsIncrementalExportResponse struct {
	TicketForms []TicketForm `json:"ticket_forms"`
	ticketingTicketFormObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
type TicketingTicketFormsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#create-ticket-form
func (s *TicketingTicketFormsService) Create(
	ctx context.Context,
	payload TicketFormPayload,
) (TicketFormResponse, error) {
	return createRequest[TicketFormResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#show-ticket-form
func (s *TicketingTicketFormsService) Show(
	ctx context.Context,
	id TicketFormID,
) (TicketFormResponse, error) {
	return showRequest[TicketFormID, TicketFormResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#list-ticket-forms
func (s *TicketingTicketFormsService) List(
	ctx context.Context,
	pageHandler func(response TicketFormsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#update-ticket-form
func (s *TicketingTicketFormsService) Update(
	ctx context.Context,
	id TicketFormID,
	payload TicketFormPayload,
) (TicketFormResponse, error) {
	return updateRequest[TicketFormID, TicketFormResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#delete-ticket-form
func (s *TicketingTicketFormsService) Delete(
	ctx context.Context,
	id TicketFormID,
) error {
	return deleteRequest[TicketFormID, TicketFormResponse](ctx, s.c, id)
}
