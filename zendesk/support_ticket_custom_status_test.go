package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportTicketCustomStatus_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_custom_status/list_page_1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_statuses",
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_custom_status/list_page_2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_statuses.json",
				Query: url.Values{
					"page": []string{"2"},
				},
			},
		),
	})

	actual := []zendesk.CustomStatus{}

	if err := z.Support().CustomStatuses().List(ctx, func(response zendesk.CustomStatusesResponse) error {
		actual = append(actual, response.CustomStatuses...)

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 5 {
		t.Fatalf("expected 5 statuses, got %d", len(actual))
	}
}

func Test_SupportTicketCustomStatus_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/ticket_custom_status/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_statuses/9001",
			},
		),
	})

	expectedCustomStatusID := zendesk.CustomStatusID(9001)

	actual, err := z.Support().CustomStatuses().Show(ctx, expectedCustomStatusID)
	if err != nil {
		t.Fatal(err)
	}

	if actual.ID != expectedCustomStatusID {
		t.Fatalf("expected ID: %d - got ID: %d", expectedCustomStatusID, actual.ID)
	}
}

func Test_SupportTicketCustomStatus_Create_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/ticket_custom_status/create_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/custom_statuses",
			},
		),
	})

	newCustomStatusLabel := "A Test Status for new tickets"

	newCustomStatus := zendesk.CustomStatusPayload{
		CustomStatus: struct {
			Active         bool   `json:"active"`
			AgentLabel     string `json:"agent_label"`
			StatusCategory string `json:"status_category"`
		}{
			Active:         true,
			AgentLabel:     newCustomStatusLabel,
			StatusCategory: zendesk.StatusNew,
		},
	}

	actual, err := z.Support().CustomStatuses().Create(ctx, newCustomStatus)
	if err != nil {
		t.Fatal(err)
	}

	if actual.CustomStatus.AgentLabel != newCustomStatusLabel {
		t.Fatalf("expected ID: %s - got ID: %s", newCustomStatusLabel, actual.CustomStatus.AgentLabel)
	}
}
