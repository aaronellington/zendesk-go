package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SuspendedTicketsResponse struct {
	SuspendedTickets []SuspendedTicket `json:"suspended_tickets"`
	CursorPaginationResponse
}

type SuspendedTicket struct {
	ID        SuspendedTicketID     `json:"id"`
	Subject   string                `json:"subject"`
	Cause     string                `json:"cause"`
	Author    SuspendedTicketAuthor `json:"author"`
	CauseID   int                   `json:"cause_id"`
	TicketID  *TicketID             `json:"ticket_id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Recipient string                `json:"recipient"`
}

type SuspendedTicketAuthor struct {
	ID    *UserID `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/
type SuspendedTicketService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/#list-suspended-tickets
func (s SuspendedTicketService) List(
	ctx context.Context,
	pageHandler func(response SuspendedTicketsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/suspended_tickets?%s", query.Encode())

	for {
		target := SuspendedTicketsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ZendeskRequest(request, &target); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/#recover-multiple-suspended-tickets
func (s *SuspendedTicketService) RecoverMultiple(ctx context.Context, ids []SuspendedTicketID) error {
	stringIDs := []string{}
	for _, id := range ids {
		stringIDs = append(stringIDs, fmt.Sprintf("%d", id))
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/suspended_tickets/recover_many?ids=%s", strings.Join(stringIDs, ",")),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(request, nil)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/#delete-suspended-ticket
func (s *SuspendedTicketService) Delete(ctx context.Context, id SuspendedTicketID) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/suspended_tickets/%d", id),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(request, nil)
}
