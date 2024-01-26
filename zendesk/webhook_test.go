package zendesk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

// https://developer.zendesk.com/documentation/webhooks/verifying/#signing-secrets-on-new-webhooks
const ZendeskTestStaticWebhookSignature string = "dGhpc19zZWNyZXRfaXNfZm9yX3Rlc3Rpbmdfb25seQ=="

func Test_WebhookVerifySignature_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{})
	recorder := httptest.NewRecorder()
	requestFile, _ := os.Open("test_requests/zendeskGroupMembershipCreate.json")
	// testRequest := httptest.NewRequest(http.MethodPost, "/webhook/zendesk/event", requestFile)

	testRequest, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/webhook/zendesk/event",
		requestFile,
	)
	testRequest.Header.Set(zendesk.WebhookHeaderSignature, ZendeskTestStaticWebhookSignature)
	testRequest.Header.Set(zendesk.WebhookHeaderSignature, "1234")

	z.Webhook().HandleWebhook(
		func(requestBody []byte) error {
			return nil
		},
		"soemthing",
	).ServeHTTP(recorder, testRequest)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal(response.StatusCode)
	}
}
