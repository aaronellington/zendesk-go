package zendesk_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/v2/zendesk"
	"github.com/aaronellington/zendesk-go/v2/zendesk/internal/testy"
)

func TestTicketingTicketsShowSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodGet,
				"https://example.zendesk.com/api/v2/tickets/2",
				"",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/show/success.json",
			),
		},
	})

	response, err := z.Ticketing().Tickets().Show(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}

	if err := testy.Assert(
		2,
		response.Ticket.ID,
	); err != nil {
		t.Fatal(err)
	}
}

func TestTicketingTicketsShowNotFound(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodGet,
				"https://example.zendesk.com/api/v2/tickets/404",
				"",
			),
			Response: createResponse(
				t,
				http.StatusNotFound,
				"internal/test_files/ticketing/tickets/show/not_found.json",
			),
		},
	})

	_, err := z.Ticketing().Tickets().Show(ctx, 404)
	if err == nil {
		t.Fatal("should have been an error")
	}

	if err := testy.Assert(err.Error(), "Zendesk API Error: RecordNotFound"); err != nil {
		t.Fatal(err)
	}
}

func TestTicketingTicketsCreateSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodPut,
				"https://example.zendesk.com/api/v2/tickets",
				"internal/test_files/ticketing/tickets/create/payload.json",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/create/success.json",
			),
		},
	})

	response, err := z.Ticketing().Tickets().Create(ctx, zendesk.TicketPayload{
		Ticket: map[string]any{
			"status": "open",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testy.Assert(
		4,
		response.Ticket.ID,
	); err != nil {
		t.Fatal(err)
	}
}

func TestTicketingTicketsUpdateSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodPost,
				"https://example.zendesk.com/api/v2/tickets/5",
				"internal/test_files/ticketing/tickets/update/payload.json",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/update/success.json",
			),
		},
	})

	response, err := z.Ticketing().Tickets().Update(ctx, 5, zendesk.TicketPayload{
		Ticket: map[string]any{
			"status": "solved",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testy.Assert(
		5,
		response.Ticket.ID,
	); err != nil {
		t.Fatal(err)
	}
}

func TestTicketingTicketsDeleteSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodDelete,
				"https://example.zendesk.com/api/v2/tickets/5",
				"",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/delete/success.json",
			),
		},
	})

	if err := z.Ticketing().Tickets().Delete(ctx, 5); err != nil {
		t.Fatal(err)
	}
}

func TestTicketingTicketsIncrementalExportSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodGet,
				"https://example.zendesk.com/api/v2/incremental/tickets?per_page=3&start_time=125",
				"",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/incremental_export/page1.json",
			),
		},
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodGet,
				"https://example.zendesk.com/api/v2/incremental/tickets?per_page=3&start_time=250",
				"",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/tickets/incremental_export/page2.json",
			),
		},
	})

	startTime := time.Unix(125, 0)
	ticketCount := 0

	if err := z.Ticketing().Tickets().IncrementalExport(
		ctx,
		startTime,
		func(tier zendesk.TicketsIncrementalExportResponse) error {
			for range tier.Tickets {
				ticketCount++
			}

			return nil
		},
		zendesk.WithTimeBasedIncrementalExportPageSize(3),
	); err != nil {
		t.Fatal(err)
	}

	if err := testy.Assert(
		3,
		ticketCount,
	); err != nil {
		t.Fatal(err)
	}
}
