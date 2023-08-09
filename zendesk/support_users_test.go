package zendesk_test

import (
	"context"
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

	if err := study.Assert(exampleUserID, actual.ID); err != nil {
		t.Fatal(err)
	}
}
