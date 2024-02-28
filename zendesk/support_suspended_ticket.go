package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SuspendedTicketResponse struct {
	SuspendedTicket SuspendedTicket `json:"suspended_ticket"`
}

type SuspendedTicketsResponse struct {
	SuspendedTickets []SuspendedTicket `json:"suspended_tickets"`
	cursorPaginationResponse
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
	client  *client
	generic genericService[
		SuspendedTicketID,
		SuspendedTicketResponse,
		SuspendedTicketsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/#list-suspended-tickets
func (s SuspendedTicketService) List(
	ctx context.Context,
	pageHandler func(response SuspendedTicketsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
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
	return s.generic.Delete(ctx, id)
}
