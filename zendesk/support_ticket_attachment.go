package zendesk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
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

type TicketAttachmentUploadResponse struct {
	Upload struct {
		Attachment  TicketAttachment   `json:"attachment"`
		Attachments []TicketAttachment `json:"attachments"`
		Token       UploadToken        `json:"token"`
	} `json:"upload"`
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

// https://developer.zendesk.com/documentation/ticketing/using-the-zendesk-api/adding-ticket-attachments-with-the-api/
func (s TicketAttachmentService) Upload(
	ctx context.Context,
	localFilePath string,
	uploadToken UploadToken,
) (TicketAttachmentUploadResponse, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return TicketAttachmentUploadResponse{}, err
	}
	defer file.Close()

	// Attempt to identify filetype by extension. If that fails, read the first 512 bytes of the file.
	contentType := mime.TypeByExtension(filepath.Ext(localFilePath))
	if contentType == "" {
		fileHeadBuffer := make([]byte, 512)
		byteCount, err := file.Read(fileHeadBuffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return TicketAttachmentUploadResponse{}, err
			}
			fileHeadBuffer = fileHeadBuffer[:byteCount]
		}

		contentType = http.DetectContentType(fileHeadBuffer)
	}

	// The content-type header has to be overridden, specify .json to account for this
	// https://developer.zendesk.com/documentation/ticketing/using-the-zendesk-api/adding-ticket-attachments-with-the-api/
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/uploads.json",
		file,
	)
	if err != nil {
		return TicketAttachmentUploadResponse{}, err
	}

	// Set the URL Parameters filename (required) and token (optional)
	queryParams := request.URL.Query()
	queryParams.Set("filename", filepath.Base(localFilePath))
	if uploadToken != "" {
		queryParams.Set("token", string(uploadToken))
	}
	request.URL.RawQuery = queryParams.Encode()

	target := TicketAttachmentUploadResponse{}

	// Set a single request override for the content type
	if err := s.client.ZendeskRequest(
		request,
		&target,
		withContentType(contentType),
	); err != nil {
		return TicketAttachmentUploadResponse{}, err
	}

	return target, nil
}
