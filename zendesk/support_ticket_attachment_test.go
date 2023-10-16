package zendesk_test

import (
	"context"
	"net/http"
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
