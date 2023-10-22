package zendesk

import (
	"context"
	"net/http"
	"time"
)

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
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
type TicketFormService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/#list-ticket-forms
func (s TicketFormService) List(
	ctx context.Context,
) (TicketFormsResponse, error) {
	target := TicketFormsResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/ticket_forms",
		http.NoBody,
	)
	if err != nil {
		return TicketFormsResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketFormsResponse{}, err
	}

	return target, nil
}
