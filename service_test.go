package zendesk_test

import (
	"testing"

	"github.com/aaronellington/zendesk-go"
	"github.com/aaronellington/zendesk-go/internal/study"
)

func createTestService(t *testing.T, queue []study.RoundTripFunc) *zendesk.Service {
	return zendesk.New(
		"example",
		zendesk.AuthenticationToken{
			Email: "example@example.com",
			Token: "abc123",
		},
		zendesk.WithRoundTripper(study.RoundTripperQueue(t, queue)),
	)
}
