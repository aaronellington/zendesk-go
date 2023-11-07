package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/
type TicketFieldService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#json-format
type TicketField struct {
	Active      bool          `json:"active"`
	ID          TicketFieldID `json:"id"`
	Description string        `json:"description"`
}

type TicketFieldResponse struct {
	TicketField TicketField `json:"ticket_field"`
}

type TicketFieldsResponse struct {
	TicketFields []TicketField `json:"ticket_fields"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#show-ticket-field
func (s TicketFieldService) Show(ctx context.Context, id TicketFieldID) (TicketField, error) {
	target := TicketFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/ticket_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return TicketField{}, err
	}

	if err := s.client.ZendeskRequest(request, &target, false); err != nil {
		return TicketField{}, err
	}

	return target.TicketField, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#list-ticket-fields
func (s TicketFieldService) List(
	ctx context.Context,
	pageHandler func(response TicketFieldsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/ticket_fields?%s",
		query.Encode(),
	)

	for {
		target := TicketFieldsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ZendeskRequest(request, &target, true); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if !target.Meta.HasMore {
			break
		}

		endpoint = target.Links.Next
	}

	return nil
}
