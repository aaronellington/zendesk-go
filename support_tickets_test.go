package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go"
	"github.com/aaronellington/zendesk-go/internal/study"
)

func Test_SupportTicketsShow_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000",
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	actual, err := z.Support().Tickets().Show(ctx, exampleTicketID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleTicketID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketsShow_401(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnauthorized,
				FilePath:   "test_files/responses/support/tickets/show_401.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000",
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	_, actual := z.Support().Tickets().Show(ctx, exampleTicketID)
	if actual == nil {
		t.Fatal("should have been an error")
	}

	expected := "Couldn't authenticate you"

	if err := study.Assert(expected, actual.Error()); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketsShow_404(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusNotFound,
				FilePath:   "test_files/responses/support/tickets/show_404.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000",
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	_, actual := z.Support().Tickets().Show(ctx, exampleTicketID)
	if actual == nil {
		t.Fatal("should have been an error")
	}

	expected := "RecordNotFound"

	if err := study.Assert(expected, actual.Error()); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketsShow_404_Wrong_Subdomain(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusNotFound,
				FilePath:   "test_files/responses/support/tickets/show_404_wrong_subdomain.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/tickets/1000",
			},
		),
	})

	var exampleTicketID zendesk.TicketID = 1000

	_, actual := z.Support().Tickets().Show(ctx, exampleTicketID)
	if actual == nil {
		t.Fatal("should have been an error")
	}

	expected := "No help desk at example.zendesk.com"

	if err := study.Assert(expected, actual.Error()); err != nil {
		t.Fatal(err)
	}
}
