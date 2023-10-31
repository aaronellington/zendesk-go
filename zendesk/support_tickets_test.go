package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
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

func Test_SupportTicketsShow_404_Wrong_SubDomain(t *testing.T) {
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

func Test_Support_Tickets_IncrementalExport(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/incremental_export_page1.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/incremental/tickets.json",
				Query: url.Values{
					"per_page":   []string{"2"},
					"start_time": []string{"0"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/incremental_export_page2.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/incremental/tickets.json",
				Query: url.Values{
					"per_page":   []string{"2"},
					"start_time": []string{"250"},
				},
			},
		),
	})

	tickets := []zendesk.Ticket{}

	if err := z.Support().Tickets().IncrementalExport(ctx, time.Unix(0, 0), 2, func(response zendesk.TicketsIncrementalExportResponse) error {
		tickets = append(tickets, response.Tickets...)

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	expectedTicketCount := 3

	if err := study.Assert(expectedTicketCount, len(tickets)); err != nil {
		t.Fatal(err)
	}
}

func Test_Support_Tickets_Merge_Success(t *testing.T) {
	ctx := context.Background()
	targetTicket := zendesk.TicketID(1000)
	sourceTickets := []zendesk.TicketID{2565, 3323}

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/merge_success.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/merge", targetTicket),
			},
		),
	})

	actual, err := z.Support().Tickets().Merge(ctx, targetTicket, zendesk.MergeRequestPayload{
		IDs:                   sourceTickets,
		SourceComment:         "test",
		SourceCommentIsPublic: true,
		TargetComment:         "done",
		TargetCommentIsPublic: true,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := "queued"

	if err := study.Assert(expected, actual.JobStatus.Status); err != nil {
		t.Fatal(err)
	}
}

func Test_Support_Tickets_Merge_TicketID_Closed(t *testing.T) {
	ctx := context.Background()
	targetTicket := zendesk.TicketID(1000)
	sourceTickets := []zendesk.TicketID{1000, 3323}

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusBadRequest,
				FilePath:   "test_files/responses/support/tickets/merge_source_invalid.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/merge", targetTicket),
			},
		),
	})

	_, err := z.Support().Tickets().Merge(ctx, targetTicket, zendesk.MergeRequestPayload{
		IDs:                   sourceTickets,
		SourceComment:         "test",
		SourceCommentIsPublic: true,
		TargetComment:         "done",
		TargetCommentIsPublic: true,
	})
	if err == nil {
		t.Fatal("should have had an error")
	}

	expected := "SourceInvalid"

	if err := study.Assert(expected, err.Error()); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTickets_ListProblemTicketIncidents_200(t *testing.T) {
	ctx := context.Background()
	problemTicketID := zendesk.TicketID(2000)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/list_problem_ticket_incidents_page_1.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/incidents.json", problemTicketID),
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/tickets/list_problem_ticket_incidents_page_2.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/tickets/%d/incidents.json", problemTicketID),
				Query: url.Values{
					"page[size]":  []string{"100"},
					"page[after]": []string{"aCursor=="},
				},
			},
		),
	})

	linkedIncidents := []zendesk.Ticket{}

	if err := z.Support().Tickets().ListProblemTicketIncidents(
		ctx,
		problemTicketID,
		func(response zendesk.ListProblemTicketIncidentsResponse) error {
			linkedIncidents = append(linkedIncidents, response.Tickets...)

			return nil
		}); err != nil {
		t.Fatal(err)
	}

	if len(linkedIncidents) != 4 {
		t.Fatalf("expected %d incidents, got %d", 4, len(linkedIncidents))
	}
}
