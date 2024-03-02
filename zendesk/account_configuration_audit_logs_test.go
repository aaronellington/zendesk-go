package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Account_Configuration_Audit_Logs_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs",
				Query: url.Values{
					"page[size]": []string{"100"},
					// "sort":       []string{"created_at"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs.json",
				Query: url.Values{
					"page[size]":  []string{"4"},
					"page[after]": []string{"eyJvIjoiLWNyZWF0ZWRfYXQsLWlkIiwidiI2IlpFZzE5MlFBQUFBQWFaRlZmdHVvRUFBQSJ9"},
					"sort":        []string{"created_at"},
				},
			},
		),
	})

	allItems := []zendesk.AuditLog{}

	if err := z.AccountConfiguration().AuditLogs().List(ctx, func(response zendesk.AuditLogsResponse) error {
		allItems = append(allItems, response.AuditLogs...)

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if len(allItems) != 8 {
		t.Fatalf("Did not get all expected results; Only have %d items in slice", len(allItems))
	}
}

func Test_Account_Configuration_Audit_Logs_List_With_Modifier_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs",
				Query: url.Values{
					"page[size]": []string{"100"},
					// "sort":             []string{"created_at"},
					// "filter[actor_id]": []string{"-1"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs.json",
				Query: url.Values{
					"page[size]":  []string{"4"},
					"page[after]": []string{"eyJvIjoiLWNyZWF0ZWRfYXQsLWlkIiwidiI2IlpFZzE5MlFBQUFBQWFaRlZmdHVvRUFBQSJ9"},
					"sort":        []string{"created_at"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs",
				Query: url.Values{
					"page[size]": []string{"100"},
					// "sort":       []string{"created_at"},
					// "filter[actor_id]": []string{"1234"},
					// "filter[action]":   []string{"create"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/audit_log/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/audit_logs.json",
				Query: url.Values{
					"page[size]":  []string{"4"},
					"page[after]": []string{"eyJvIjoiLWNyZWF0ZWRfYXQsLWlkIiwidiI2IlpFZzE5MlFBQUFBQWFaRlZmdHVvRUFBQSJ9"},
					"sort":        []string{"created_at"},
				},
			},
		),
	})

	// Single modifier
	if err := z.AccountConfiguration().AuditLogs().List(
		ctx,
		func(response zendesk.AuditLogsResponse) error {
			return nil
		},
		zendesk.WithFilterForActorID(zendesk.ActorID(-1)),
	); err != nil {
		t.Fatal(err)
	}

	// 2 modifiers
	if err := z.AccountConfiguration().AuditLogs().List(
		ctx,
		func(response zendesk.AuditLogsResponse) error {
			return nil
		},
		zendesk.WithFilterForActorID(zendesk.ActorID(1234)),
		zendesk.WithFilterForAction(zendesk.Create),
	); err != nil {
		t.Fatal(err)
	}
}
