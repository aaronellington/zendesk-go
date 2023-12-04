package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Support_Ticket_Comments_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_comments/list_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000/comments",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	actual := []zendesk.TicketComment{}

	if err := z.Support().TicketComments().ListByTicketID(
		ctx,
		exampleTicketID,
		func(response zendesk.TicketCommentResponse) error {
			actual = append(actual, response.Comments...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(actual))
	}
}

func Test_Support_Ticket_Comments_List_SystemComment_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_comments/show_systemComment_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000/comments",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	actual := []zendesk.TicketComment{}

	if err := z.Support().TicketComments().ListByTicketID(
		ctx,
		exampleTicketID,
		func(response zendesk.TicketCommentResponse) error {
			actual = append(actual, response.Comments...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(actual))
	}
}
