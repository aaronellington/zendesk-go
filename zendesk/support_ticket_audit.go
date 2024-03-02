package zendesk

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type TicketAudit struct {
	ID        TicketAuditID      `json:"id"`
	TicketID  int                `json:"ticket_id"`
	CreatedAt time.Time          `json:"created_at"`
	AuthorID  int                `json:"author_id"`
	Events    []TicketAuditEvent `json:"events"`
}

type TicketAuditEvent struct {
	ID            TicketAuditEventID `json:"id"`
	Type          string             `json:"type"`
	FieldName     string             `json:"field_name"`
	PreviousValue any                `json:"previous_value"`
	Value         any                `json:"value"`
}

type TicketAuditsResponse struct {
	Audits []TicketAudit `json:"audits"`
	cursorPaginationResponse
}
type TicketAuditResponse struct {
	Audit TicketAudit `json:"audit"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_audits/
type TicketAuditService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_audits/#list-audits-for-a-ticket
func (s TicketAuditService) ListForTicket(
	ctx context.Context,
	ticketID TicketID,
	pageHandler func(response TicketAuditsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/tickets/%d/audits?%s", ticketID, query.Encode())

	return genericList(
		ctx,
		s.client,
		endpoint,
		pageHandler,
	)
}
