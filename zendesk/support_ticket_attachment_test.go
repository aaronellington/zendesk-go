package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportTicketAttachmentShow_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_attachment/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/attachments/1001",
			},
		),
	})

	var exampleAttachmentID zendesk.AttachmentID = 1001

	actual, err := z.Support().TicketAttachments().Show(ctx, exampleAttachmentID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleAttachmentID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketAttachmentUpload_png_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_attachment/upload_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/uploads.json",
				Query: url.Values{
					"filename": []string{"gopher.png"},
				},
				Validator: func(r *http.Request) error {
					expectedContentType := "image/png"
					if r.Header.Get("Content-Type") != expectedContentType {
						return fmt.Errorf(
							"expected content type to be '%s', got '%s'",
							expectedContentType,
							r.Header.Get("Content-Type"),
						)
					}
					return nil
				},
			},
		),
	})

	expectedUploadToken := zendesk.UploadToken("FakeUploadToken")

	actual, err := z.Support().TicketAttachments().Upload(ctx, "test_files/responses/support/ticket_attachment/attachments/gopher.png", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedUploadToken, actual.Upload.Token); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketAttachmentUpload_svg_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_attachment/upload_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/uploads.json",
				Query: url.Values{
					"filename": []string{"gopher.svg"},
				},
				Validator: func(r *http.Request) error {
					expectedContentType := "image/svg+xml"
					if r.Header.Get("Content-Type") != expectedContentType {
						return fmt.Errorf(
							"expected content type to be '%s', got '%s'",
							expectedContentType,
							r.Header.Get("Content-Type"),
						)
					}
					return nil
				},
			},
		),
	})

	expectedUploadToken := zendesk.UploadToken("FakeUploadToken")

	actual, err := z.Support().TicketAttachments().Upload(ctx, "test_files/responses/support/ticket_attachment/attachments/gopher.svg", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedUploadToken, actual.Upload.Token); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketAttachmentUpload_noFileExt_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_attachment/upload_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/uploads.json",
				Query: url.Values{
					"filename": []string{"gopherNoFileExt"},
				},
				Validator: func(r *http.Request) error {
					expectedContentType := "image/png"
					actualContentType := r.Header.Get("Content-Type")
					return study.Assert(expectedContentType, actualContentType)
				},
			},
		),
	})

	expectedUploadToken := zendesk.UploadToken("FakeUploadToken")

	actual, err := z.Support().TicketAttachments().Upload(ctx, "test_files/responses/support/ticket_attachment/attachments/gopherNoFileExt", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedUploadToken, actual.Upload.Token); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketAttachmentUpload_Multiple_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_attachment/upload_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/uploads.json",
				Query: url.Values{
					"filename": []string{"gopher.png"},
				},
				Validator: func(r *http.Request) error {
					expectedContentType := "image/png"
					actualContentType := r.Header.Get("Content-Type")
					return study.Assert(expectedContentType, actualContentType)
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_attachment/upload_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/uploads.json",
				Query: url.Values{
					"filename": []string{"gopher.svg"},
					"token":    []string{"FakeUploadToken"},
				},
				Validator: func(r *http.Request) error {
					expectedContentType := "image/svg+xml"
					actualContentType := r.Header.Get("Content-Type")
					return study.Assert(expectedContentType, actualContentType)
				},
			},
		),
	})

	upload1, err := z.Support().TicketAttachments().Upload(ctx, "test_files/responses/support/ticket_attachment/attachments/gopher.png", "")
	if err != nil {
		t.Fatal(err)
	}

	upload2, err := z.Support().TicketAttachments().Upload(ctx, "test_files/responses/support/ticket_attachment/attachments/gopher.svg", upload1.Upload.Token)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(upload2.Upload.Token, upload1.Upload.Token); err != nil {
		t.Fatal(err)
	}
}
