package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_SupportUsersShow_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/users/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/users/1000",
			},
		),
	})

	var exampleUserID zendesk.UserID = 1000

	actual, err := z.Support().Users().Show(ctx, exampleUserID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(exampleUserID, actual.User.ID); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportUsersCreate_201(t *testing.T) {
	ctx := context.Background()
	email := "kren+newEmail@chandrila.com"
	name := "Kylo Ren"

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusCreated,
				FilePath:   "test_files/responses/support/users/create_201.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/users",
			},
		),
	})

	actual, err := z.Support().Users().Create(ctx, zendesk.UserPayload{
		User: struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			Email: email,
			Name:  name,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(actual.User.Name, name); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportUsersCreate_422(t *testing.T) {
	ctx := context.Background()
	email := "kren+newEmailTaken@chandrila.com"
	name := ""

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusUnprocessableEntity,
				FilePath:   "test_files/responses/support/users/create_422.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodPost,
				Path:   "/api/v2/users",
			},
		),
	})

	_, err := z.Support().Users().Create(ctx, zendesk.UserPayload{
		User: struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			Email: email,
			Name:  name,
		},
	})
	if err == nil {
		t.Fatal("expected error for this test")
	}

	zdErr, ok := err.(*zendesk.Error)
	if !ok {
		t.Fatal("expected zendesk err")
	}

	expectedErrorString := "RecordInvalid. Error details: [name: ValueTooShort - Name: is too short (minimum one character)]"
	if err := study.Assert(expectedErrorString, zdErr.Error()); err != nil {
		t.Fatal(err)
	}
}

func Test_SupportUsersShowWithSideloads_429Retry(t *testing.T) {
	ctx := context.Background()
	userID := zendesk.UserID(1000)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusTooManyRequests,
				FilePath:   "test_files/responses/errors/api_rate_limit_exceeded_429_inconsistent.txt",
				ResponseModifiers: []study.ResponseModifier{
					study.WithResponseHeaders(map[string][]string{
						"retry-after":  {"0"},
						"Content-Type": {""},
					}),
				},
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/users/%d", userID),
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/users/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/users/%d", userID),
			},
		),
	})

	actual, err := z.Support().Users().Show(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(actual.User.ID, userID); err != nil {
		t.Fatal(err)
	}
}
