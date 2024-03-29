package zendesk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/#json-format
type TicketAttachment struct {
	ContentType           string            `json:"content_type"`
	ContentURL            string            `json:"content_url"`
	Deleted               bool              `json:"deleted"`
	FileName              string            `json:"file_name"`
	Height                uint64            `json:"height"`
	ID                    AttachmentID      `json:"id"`
	Inline                bool              `json:"inline"`
	MalwareAccessOverride bool              `json:"malware_access_override"`
	MalwareScanResult     MalwareScanResult `json:"malware_scan_result"`
	MappedContentURL      string            `json:"mapped_content_url"`
	Size                  uint64            `json:"size"`
	URL                   string            `json:"url"`
	Width                 uint64            `json:"width"`
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

/*
When working with attachments, we must ensure that the domain that the content is hosted on
is *.zendesk.com before sending authentication credentials.

- https://developer.zendesk.com/documentation/ticketing/managing-tickets/working-with-url-properties/

- https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/#json-format
*/
func (s TicketAttachmentService) Download(
	ctx context.Context,
	contentURL string,
	writer io.Writer,
) error {
	var response *http.Response

	var err error

	if strings.Contains(contentURL, "zendesk.com") {
		response, err = s.client.ZendeskGetRequest(ctx, contentURL)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	} else {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, contentURL, nil)
		if err != nil {
			return err
		}

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}

	_, err = io.Copy(writer, response.Body)

	return err
}

func (s *TicketAttachmentService) DownloadToFile(
	ctx context.Context,
	contentURL string,
	filePath string,
) error {
	outfile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	return s.Download(ctx, contentURL, outfile)
}

// https://developer.zendesk.com/documentation/ticketing/using-the-zendesk-api/adding-ticket-attachments-with-the-api/
func (s TicketAttachmentService) Upload(
	ctx context.Context,
	localFilePath string,
	uploadToken UploadToken,
) (TicketAttachmentUploadResponse, error) {
	return s.UploadWithFilename(
		ctx,
		localFilePath,
		filepath.Base(localFilePath),
		uploadToken,
	)
}

func (s TicketAttachmentService) UploadWithFilename(
	ctx context.Context,
	localFilePath string,
	filename string,
	uploadToken UploadToken,
) (TicketAttachmentUploadResponse, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return TicketAttachmentUploadResponse{}, err
	}
	defer file.Close()

	buf := &bytes.Buffer{}

	_, err = buf.ReadFrom(file)
	if err != nil {
		return TicketAttachmentUploadResponse{}, err
	}

	// The content-type header has to be overridden, specify .json to account for this
	// https://developer.zendesk.com/documentation/ticketing/using-the-zendesk-api/adding-ticket-attachments-with-the-api/
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/uploads.json",
		buf,
	)
	if err != nil {
		return TicketAttachmentUploadResponse{}, err
	}

	// Set the content-length header
	request.Header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	contentType := mime.TypeByExtension(filepath.Ext(localFilePath))
	if contentType == "" {
		contentType = http.DetectContentType(buf.Bytes())
	}

	// Set the content-type header
	request.Header.Set("Content-Type", contentType)

	// Set the URL Parameters filename (required) and token (optional)
	queryParams := request.URL.Query()
	queryParams.Set("filename", filepath.Base(filename))

	if uploadToken != "" {
		queryParams.Set("token", string(uploadToken))
	}

	request.URL.RawQuery = queryParams.Encode()

	target := TicketAttachmentUploadResponse{}

	if err := s.client.ZendeskRequest(
		request,
		&target,
	); err != nil {
		return TicketAttachmentUploadResponse{}, err
	}

	return target, nil
}
