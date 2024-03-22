package zendesk_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	z, _ := createTestRealTimeChatWebsocketService(
		t,
		ctx,
		settings{},
	)

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				testError <- err
				return
			}
		}
	}()

	timeout := time.NewTimer(time.Second * 5)
	interval := time.NewTicker(time.Millisecond * 250)
	for {
		select {
		case err := <-testError:
			if err != nil {
				t.Fatal(err)
			}

		case <-timeout.C:
			t.Fatal("did not record a connection within timeout")

		case <-interval.C:
			startedTime := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetConnectionStartedTime()
			if startedTime != nil {
				return
			}
		}
	}
}

func TestService_SessionExpiring(t *testing.T) {
}

func TestService_PongFailure(t *testing.T) {
}

func TestService_HandleRetryableError(t *testing.T) {
}

func TestService_HandleFatalError(t *testing.T) {
}

func TestService_BadCredentials(t *testing.T) {
}

func TestService_ConfirmNoStaleGoRoutines(t *testing.T) {
	// Do this at the end of every test.
}
