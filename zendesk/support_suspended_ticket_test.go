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

func Test_SupportSuspendedTicket_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/suspended_ticket/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/suspended_tickets/12345",
			},
		),
	})

	var exampleSuspendedTicketID zendesk.SuspendedTicketID = 12345

	actual, err := z.Support().SuspendedTickets().Show(ctx, exampleSuspendedTicketID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleSuspendedTicketID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportSuspendedTicket_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/suspended_ticket/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/suspended_tickets",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/suspended_ticket/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/suspended_tickets.json",
				Query: url.Values{
					"page[size]":  []string{"1"},
					"page[after]": []string{"aCursor"},
				},
			},
		),
	})

	expectedSuspendedTicketsLen := 2
	actualSuspendedTicketsLen := 0

	if err := z.Support().SuspendedTickets().List(ctx,
		func(response zendesk.SuspendedTicketsResponse) error {
			for range response.SuspendedTickets {
				actualSuspendedTicketsLen++
			}

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedSuspendedTicketsLen, actualSuspendedTicketsLen); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportSuspendedTicket_Delete_200(t *testing.T) {
	ctx := context.Background()

	var exampleSuspendedTicketID zendesk.SuspendedTicketID = 12345

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseNoContent{
				StatusCode: http.StatusNoContent,
			},
			study.ExpectedTestRequest{
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/api/v2/suspended_tickets/%d", exampleSuspendedTicketID),
			},
		),
	})

	if err := z.Support().SuspendedTickets().Delete(ctx, exampleSuspendedTicketID); err != nil {
		t.Fatal(err)
	}
}
