package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestChatOAuthClient_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/chat/oauth_client/list_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/oauth/clients",
			},
		),
	})

	actual, err := z.LiveChat().Chat().OAuthClientService().List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 2 {
		t.Fatalf("got %d clients, expected 2 clients", len(actual))
	}
}
