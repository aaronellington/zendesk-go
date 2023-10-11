package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/#json-format
type TicketAttachment struct {
	ContentType string       `json:"content_type"`
	ContentURL  string       `json:"content_url"`
	FileName    string       `json:"file_name"`
	ID          AttachmentID `json:"id"`
	Size        uint64       `json:"size"`
}

type TicketAttachmentResponse struct {
	Attachment TicketAttachment `json:"attachment"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/
type TicketAttachmentService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/#show-attachment
func (s TicketAttachmentService) Show(
	ctx context.Context,
	attachmentID AttachmentID,
) (TicketAttachment, error) {
	target := TicketAttachmentResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/attachments/%d", attachmentID),
		http.NoBody,
	)
	if err != nil {
		return TicketAttachment{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketAttachment{}, err
	}

	return target.Attachment, nil
}
