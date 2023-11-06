package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Client_429(t *testing.T) {
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
				StatusCode: http.StatusTooManyRequests,
				FilePath:   "test_files/responses/errors/api_rate_limit_exceeded_429.json",
				ResponseModifiers: []study.ResponseModifier{
					study.WithResponseHeaders(
						map[string][]string{
							"retry-after": {"0"},
						},
					),
				},
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
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusTooManyRequests,
				FilePath:   "test_files/responses/errors/api_rate_limit_exceeded_429.json",
				ResponseModifiers: []study.ResponseModifier{
					study.WithResponseHeaders(
						map[string][]string{
							"retry-after": {"0"},
						},
					),
				},
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
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusTooManyRequests,
				FilePath:   "test_files/responses/errors/api_rate_limit_exceeded_429.json",
				ResponseModifiers: []study.ResponseModifier{
					study.WithResponseHeaders(
						map[string][]string{
							"retry-after": {"0"},
						},
					),
				},
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
