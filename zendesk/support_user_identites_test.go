package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportUserIdentitiesList_200(t *testing.T) {
	ctx := context.Background()
	var userID zendesk.UserID = 1000

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/user_identities/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/users/%d/identities", userID),
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/user_identities/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/users/%d/identities.json", userID),
				Query: url.Values{
					"page[size]":  []string{"100"},
					"page[after]": []string{"aCursor="},
				},
			},
		),
	})

	actualIdentities := []zendesk.UserIdentity{}

	if err := z.Support().UserIdentities().List(
		ctx,
		userID,
		func(response zendesk.UserIdentitiesResponse) error {
			actualIdentities = append(actualIdentities, response.Identities...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	expectedNumberOfIdentities := 3

	if err := study.Assert(expectedNumberOfIdentities, len(actualIdentities)); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportTicketFormsCreate_201(t *testing.T) {
	ctx := context.Background()
	var userID zendesk.UserID = 1000
	userEmail := "kren+newEmail@chandrila.com"

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/user_identities/create_201.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/api/v2/users/%d/identities", userID),
			},
		),
	})

	actual, err := z.Support().UserIdentities().Create(
		ctx,
		userID, zendesk.UserIdentityPayload{
			Identity: struct {
				Type            string `json:"type"`
				Value           string `json:"value"`
				SkipVerifyEmail bool   `json:"skip_verify_email"`
			}{
				Type:            "email",
				Value:           userEmail,
				SkipVerifyEmail: true,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(userEmail, actual.Identity.Value); err != nil {
		t.Fatal(err)
	}
}
