package zendesk_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Client_422(t *testing.T) {
	ctx := context.Background()
	closedTicketID := zendesk.TicketID(1000)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnprocessableEntity,
				FilePath:   "test_files/responses/errors/422_record_invalid.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/api/v2/tickets/%d", closedTicketID),
			},
		),
	})

	_, err := z.Support().Tickets().Update(ctx, closedTicketID, zendesk.TicketPayload{
		Ticket: zendesk.TicketComment{
			Body: "This is a test comment, which is being added to a closed ticket.",
		},
	})
	if err == nil {
		t.Fatal("expected an error - did not receive one")
	}

	zendeskGoError := &zendesk.Error{}
	isZendeskGoError := errors.As(err, &zendeskGoError)

	if !isZendeskGoError {
		t.Fatalf("expected a custom zendesk-go error, got: %v", err)
	}

	// Check to confirm that we got a 422 error
	if !zendeskGoError.ImmutableRecord() {
		t.Fatal("did not receive an immutable error")
	}
}

func Test_Client_429(t *testing.T) {
	ctx := context.Background()
	allRequestsMade := false

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
				Validator: func(r *http.Request) error {
					allRequestsMade = true

					return nil
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

	if !allRequestsMade {
		t.Fatal("not all requests were made")
	}
}

func Test_Client_429_Retries_Exceeded(t *testing.T) {
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
	})

	tickets := []zendesk.Ticket{}

	err := z.Support().Tickets().IncrementalExport(ctx, time.Unix(0, 0), 2, func(response zendesk.TicketsIncrementalExportResponse) error {
		tickets = append(tickets, response.Tickets...)

		return nil
	})
	if err == nil {
		t.Fatalf("expected to get error")
	}

	zdErr, ok := err.(*zendesk.Error)
	if !ok {
		t.Fatalf("expected to get error of type zendesk.Error, received: %T", err)
	}

	if zdErr.Response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected to get HTTP 429, got: %d", zdErr.Response.StatusCode)
	}
}

func Test_Client_HTML_Error_Received(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusInternalServerError,
				FilePath:   "test_files/responses/errors/html_error_response.html",
				ResponseModifiers: []study.ResponseModifier{
					study.WithResponseHeaders(
						map[string][]string{
							"Content-Type": {""},
						},
					),
				},
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
	})

	tickets := []zendesk.Ticket{}

	err := z.Support().Tickets().IncrementalExport(ctx, time.Unix(0, 0), 2, func(response zendesk.TicketsIncrementalExportResponse) error {
		tickets = append(tickets, response.Tickets...)

		return nil
	})
	if err == nil {
		t.Fatalf("expected to get error")
	}

	zdErr, ok := err.(*zendesk.Error)
	if !ok {
		t.Fatalf("expected to get error of type zendesk.Error, received: %T", err)
	}

	if zdErr.Response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected to get 500 error with HTML response body, got: %d", zdErr.Response.StatusCode)
	}
}

func Test_Client_ECONNRESET_Retry_Success(t *testing.T) {
	ctx := context.Background()
	allRequestsMade := false

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseURLError{
				URLError: &url.Error{
					Op: "Get",
					Err: &net.OpError{
						Op:  "accept",
						Net: "tcp",
						Err: syscall.ECONNRESET,
					},
				},
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
					"start_time": []string{"0"},
				},
				Validator: func(r *http.Request) error {
					allRequestsMade = true

					return nil
				},
			},
		),
	})

	tickets := []zendesk.Ticket{}

	err := z.Support().Tickets().IncrementalExport(ctx, time.Unix(0, 0), 2, func(response zendesk.TicketsIncrementalExportResponse) error {
		tickets = append(tickets, response.Tickets...)

		return nil
	})
	if err != nil {
		t.Fatalf("expected to get error")
	}

	if !allRequestsMade {
		t.Fatal("expected to retry on temporary error")
	}
}

func Test_Client_ECONNRESET_Retry_Exhausted(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseURLError{
				URLError: &url.Error{
					Op: "Get",
					Err: &net.OpError{
						Op:  "accept",
						Net: "tcp",
						Err: syscall.ECONNRESET,
					},
				},
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
			&study.TestResponseURLError{
				URLError: &url.Error{
					Op: "Get",
					Err: &net.OpError{
						Op:  "accept",
						Net: "tcp",
						Err: syscall.ECONNRESET,
					},
				},
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
			&study.TestResponseURLError{
				URLError: &url.Error{
					Op: "Get",
					Err: &net.OpError{
						Op:  "accept",
						Net: "tcp",
						Err: syscall.ECONNRESET,
					},
				},
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/incremental/tickets.json",
				Query: url.Values{
					"per_page":   []string{"2"},
					"start_time": []string{"0"},
				},
				Validator: func(r *http.Request) error {
					if r.Header.Get(zendesk.RequestHeaderRetryAttempts) != "3" {
						t.Fatalf("expected three attempts, got %s", r.Header.Get(zendesk.RequestHeaderRetryAttempts))
					}

					return nil
				},
			},
		),
	})

	tickets := []zendesk.Ticket{}

	err := z.Support().Tickets().IncrementalExport(ctx, time.Unix(0, 0), 2, func(response zendesk.TicketsIncrementalExportResponse) error {
		tickets = append(tickets, response.Tickets...)

		return nil
	})
	if err == nil {
		t.Fatalf("expected to get error")
	}

	networkErr, ok := err.(*url.Error)
	if !ok {
		t.Fatalf("expected network error, got error: %s", err.Error())
	}

	if !errors.Is(err, syscall.ECONNRESET) {
		t.Fatalf("did not get correct network error, got: %s", networkErr.Err.Error())
	}
}
