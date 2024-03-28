package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/v2/zendesk/internal/testy"
)

func TestTicketingUsersShowSuccess(t *testing.T) {
	ctx := context.Background()

	z := getTestInstance(t, []http.RoundTripper{
		testy.RequestResponseTester{
			Request: createBaseRequest(
				t,
				ctx,
				http.MethodGet,
				"https://example.zendesk.com/api/v2/users/2",
				"",
			),
			Response: createResponse(
				t,
				http.StatusOK,
				"internal/test_files/ticketing/users/show/success.json",
			),
		},
	})

	response, err := z.Ticketing().Users().Show(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}

	if err := testy.Assert(
		2,
		response.User.ID,
	); err != nil {
		t.Fatal(err)
	}
}
