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

const ZendeskTestStaticWebhookSignature string = "dGhpc19zZWNyZXRfaXNfZm9yX3Rlc3Rpbmdfb25seQ=="

func Test_WebhookVerifySignature_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{})
	recorder := httptest.NewRecorder()
	requestFile, _ := os.Open("test_files/requests/webhook/user/group_membership_created.json")

	testRequest, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/webhook/zendesk/event",
		requestFile,
	)
	testRequest.Header.Set(zendesk.WebhookHeaderSignature, "10IqYzYTLHRftNsNE+im0DeOM6/JactIRuy0XCHJ9B8=")
	testRequest.Header.Set(zendesk.WebhookHeaderSignatureTimestamp, "1234")

	z.Webhook().HandleWebhookEvent(
		zendesk.WebhookEventHandlers{
			WebhookEventUserGroupMembershipCreated: func(eventData zendesk.WebhookEventUserGroupMembershipCreatedPayload) error {
				return nil
			},
		},
		zendesk.WithSigningSecret(ZendeskTestStaticWebhookSignature),
	).ServeHTTP(recorder, testRequest)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal(response.StatusCode)
	}
}

func Test_WebhookSkipVerifySignature_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{})
	recorder := httptest.NewRecorder()
	requestFile, _ := os.Open("test_files/requests/webhook/user/group_membership_created.json")

	testRequest, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/webhook/zendesk/event",
		requestFile,
	)
	testRequest.Header.Set(zendesk.WebhookHeaderSignature, "10IqYzYTLHRftNsNE+im0DeOM6/JactIRuy0XCHJ9B8=")
	testRequest.Header.Set(zendesk.WebhookHeaderSignatureTimestamp, "1234")

	z.Webhook().HandleWebhookEvent(
		zendesk.WebhookEventHandlers{
			WebhookEventUserGroupMembershipCreated: func(eventData zendesk.WebhookEventUserGroupMembershipCreatedPayload) error {
				return nil
			},
		},
	).ServeHTTP(recorder, testRequest)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal(response.StatusCode)
	}
}
