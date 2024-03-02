package zendesk_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_AccountConfiguration_TicketCustomRoles_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/custom_role/list_page_1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_roles",
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/custom_role/list_page_2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_roles.json",
				Query: url.Values{
					"page": []string{"2"},
				},
			},
		),
	})

	actual := []zendesk.CustomRole{}

	if err := z.AccountConfiguration().CustomRoles().List(ctx, func(response zendesk.CustomRolesResponse) error {
		actual = append(actual, response.CustomRoles...)

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 4 {
		t.Fatalf("expected 4 statuses, got %d", len(actual))
	}
}

func Test_SupportTicketCustomRole_Show_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/account_configuration/custom_role/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/custom_roles/9999",
			},
		),
	})

	expectedCustomRoleID := zendesk.CustomRoleID(9999)

	actual, err := z.AccountConfiguration().CustomRoles().Show(ctx, expectedCustomRoleID)
	if err != nil {
		t.Fatal(err)
	}

	if actual.CustomRole.ID != expectedCustomRoleID {
		t.Fatalf("expected ID: %d - got ID: %d", expectedCustomRoleID, actual.CustomRole.ID)
	}
}
