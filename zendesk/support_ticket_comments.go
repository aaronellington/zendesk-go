package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/
type TicketCommentService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#json-format
type TicketComment struct {
	ID          TicketCommentID    `json:"id"`
	Attachments []TicketAttachment `json:"attachments"`
	AuditID     AuditID            `json:"audit_id"`
	AuthorID    UserID             `json:"author_id"`
	Body        string             `json:"body"`
}

type TicketCommentResponse struct {
	Comments []TicketComment `json:"comments"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func (s TicketCommentService) ListByTicketID(
	ctx context.Context,
	ticketID TicketID,
	pageHandler func(response TicketCommentResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/tickets/%d/comments?%s",
		ticketID,
		query.Encode(),
	)

	for {
		target := TicketCommentResponse{}

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
