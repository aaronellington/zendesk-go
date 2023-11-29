package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/
type TicketCommentService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#json-format
type TicketComment struct {
	ID          TicketCommentID       `json:"id"`
	Attachments []TicketAttachment    `json:"attachments"`
	AuditID     AuditID               `json:"audit_id"`
	AuthorID    UserID                `json:"author_id"`
	Body        string                `json:"body"`
	CreatedAt   time.Time             `json:"created_at"`
	HTMLBody    string                `json:"html_body"`
	Metadata    TicketCommentMetadata `json:"metadata"`
	PlainBody   string                `json:"plain_body"`
	Public      bool                  `json:"public"`
	Type        string                `json:"type"`
	Uploads     []UploadToken         `json:"uploads"`
}

type TicketCommentMetadata struct {
	System TicketCommentMetadataSystem `json:"system"`
	Via    TicketVia                   `json:"via"`
	Flags  []uint                      `json:"flags"` // https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#comment-flags
}

type TicketCommentMetadataSystem struct {
	Client    string  `json:"client"`
	IPAddress string  `json:"ip_address"`
	Latitude  float64 `json:"latitude"`
	Location  string  `json:"location"`
	Longitude float64 `json:"longitude"`
}

type TicketCommentResponse struct {
	Comments []TicketComment `json:"comments"`
	Users    []User          `json:"users"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/#list-comments
func (s TicketCommentService) ListByTicketID(
	ctx context.Context,
	ticketID TicketID,
	pageHandler func(response TicketCommentResponse) error,
) error {
	return s.ListByTicketIDWithSideload(ctx, ticketID, nil, pageHandler)
}

func (s TicketCommentService) ListByTicketIDWithSideload(
	ctx context.Context,
	ticketID TicketID,
	sideloads []TicketCommentSideload,
	pageHandler func(response TicketCommentResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")

	if len(sideloads) > 0 {
		sideload, sideloads := string(sideloads[0]), sideloads[1:]

		for _, s := range sideloads {
			sideload = fmt.Sprintf("%s,%s", sideload, string(s))
		}

		query.Set("include", sideload)
	}

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
