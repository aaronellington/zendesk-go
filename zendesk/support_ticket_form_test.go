package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportTicketFormsShow_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_form/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/ticket_forms/123123123",
			},
		),
	})

	var exampleFormID zendesk.TicketFormID = 123123123

	actual, err := z.Support().TicketForms().Show(ctx, exampleFormID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleFormID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketFormsList_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_form/list_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/ticket_forms",
			},
		),
	})

	actual, err := z.Support().TicketForms().List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 3 {
		t.Fatal("expected 3 form objects to be returned")
	}
}
