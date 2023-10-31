package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportTicketField_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_field/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/ticket_fields/9001",
			},
		),
	})

	var exampleTicketFieldID zendesk.TicketFieldID = 9001

	actual, err := z.Support().TicketFields().Show(ctx, exampleTicketFieldID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleTicketFieldID, actual.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketField_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_field/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/ticket_fields",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_field/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/ticket_fields.json",
				Query: url.Values{
					"page[size]":  []string{"2"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	expectedFieldsLen := 4
	actualFieldsLen := 0

	if err := z.Support().TicketFields().List(ctx,
		func(response zendesk.TicketFieldsResponse) error {
			for range response.TicketFields {
				actualFieldsLen++
			}

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedFieldsLen, actualFieldsLen); err != nil {
		t.Fatal(err)
	}
}
