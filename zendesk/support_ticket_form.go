package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
type TicketFormService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#json-format
type TicketForm struct {
	Active      bool         `json:"active"`
	ID          TicketFormID `json:"id"`
	DisplayName string       `json:"display_name"`
}

type TicketFormsResponse struct {
	TicketForms []TicketForm `json:"ticket_forms"`
}

type TicketFormResponse struct {
	TicketForm TicketForm `json:"ticket_form"`
}

/*
https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#list-ticket-forms

Does not support pagination
*/
func (s TicketFormService) List(
	ctx context.Context,
) ([]TicketForm, error) {
	target := TicketFormsResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/ticket_forms",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return nil, err
	}

	return target.TicketForms, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#show-ticket-form
func (s TicketFormService) Show(
	ctx context.Context,
	ticketFormID TicketFormID,
) (TicketForm, error) {
	target := TicketFormResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/ticket_forms/%d", ticketFormID),
		http.NoBody,
	)
	if err != nil {
		return TicketForm{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketForm{}, err
	}

	return target.TicketForm, nil
}
