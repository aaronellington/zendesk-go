package zendesk_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func createTestService(t *testing.T, queue []study.RoundTripFunc) *zendesk.Service {
	return zendesk.NewService(
		"example",
		zendesk.AuthenticationToken{
			Email: "example@example.com",
			Token: "abc123",
		},
		zendesk.ChatCredentials{
			ClientID:     "fake-client-id",
			ClientSecret: "fake-client-secret",
		},
		zendesk.WithRoundTripper(study.RoundTripperQueue(t, queue)),
	)
}

func createTestRealTimeChatWebsocketService(t *testing.T, queue []study.RoundTripFunc, wsHost string) *zendesk.Service {
	return zendesk.NewService(
		"example",
		zendesk.AuthenticationToken{
			Email: "example@example.com",
			Token: "abc123",
		},
		zendesk.ChatCredentials{
			ClientID:     "fake-client-id",
			ClientSecret: "fake-client-secret",
		},
		zendesk.WithRoundTripper(study.RoundTripperQueue(t, queue)),
		zendesk.SetRealTimeChatWebsocketHost(wsHost),
	)
}

func createSuccessfulChatAuth(t *testing.T) study.RoundTripFunc {
	return study.ServeAndValidate(
		t,
		&study.TestResponseFile{
			StatusCode: 200,
			FilePath:   "test_files/responses/livechat/oauth_token_200.json",
		},
		study.ExpectedTestRequest{
			Method: http.MethodPost,
			Path:   "/oauth2/token",
			Validator: func(r *http.Request) error {
				expectedContentType := "application/x-www-form-urlencoded"
				if r.Header.Get("Content-Type") != expectedContentType {
					return fmt.Errorf(
						"expected content type to be '%s', got '%s'",
						expectedContentType,
						r.Header.Get("Content-Type"),
					)
				}

				requestBody, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}

				expectedBodyFile, err := os.Open("test_files/requests/livechat/oauth_request_body.txt")
				if err != nil {
					t.Fatal(err)
				}

				expectedBody, err := io.ReadAll(expectedBodyFile)
				if err != nil {
					t.Fatal(err)
				}

				if err := study.Assert(expectedBody, requestBody); err != nil {
					return fmt.Errorf(
						"expected body and actual body do not match - err: %s",
						err.Error(),
					)
				}

				return nil
			},
		},
	)
}
