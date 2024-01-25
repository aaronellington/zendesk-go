package zendesk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_WebhookVerification_200(t *testing.T) {
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

	actualEventData := zendesk.WebhookEventUser{}

	recorder := httptest.NewRecorder
	handler := z.Webhook().HandleWebhookUserEvent(
		"secret",
		func(e zendesk.WebhookEventUser) error {
			if e.Detail.Email != "cgibson@datto.com" {
				t.Fatal("did not match")
				return
			}

		},
	)

}
