package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

type UserID int

func Test_DeleteSession(t *testing.T) {
	ctx := context.Background()

	userID := zendesk.UserID(222)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusNoContent,
				FilePath:   "test_files/responses/account_configuration/delete_sessions/delete_user_session.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/api/v2/users/%d/sessions", userID),
			},
		),
	})

	err := z.AccountConfiguration().Sessions().BulkDelete(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
}
