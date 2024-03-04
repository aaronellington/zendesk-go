package zendesk_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	mockRealTimeChatWebsocketServer := newMockRealTimeChatWebsocketServer(testError, nil, nil, nil)

	mockServer := mockRealTimeChatWebsocketServer.createDefaultServer()

	mockServer.Start()
	defer mockServer.Close()

	rtcWSHost := strings.TrimPrefix(mockServer.URL, "http")

	z := createTestRealTimeChatWebsocketService(
		t,
		[]study.RoundTripFunc{
			createSuccessfulChatAuth(t),
		},
		fmt.Sprintf("ws%s", rtcWSHost),
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
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a connection within timeout")
	case successfulConnection := <-mockRealTimeChatWebsocketServer.state.successfulConnection:
		if !successfulConnection {
			t.Fatal("did not connect successfully")
		}

		return
	}
}

func TestRealTimeChatWebsocketStreaming_Connect_401(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	// This is our test mockserver
	mockableZendeskRTCWebsocketServer := &mockRealTimeChatWebsocketServer{
		state: state{
			successfulConnection: make(chan bool),
		},
		settings: settings{
			ValidOAuthToken: "No valid token",
		},
		testError: testError,
	}

	mockServer := mockableZendeskRTCWebsocketServer.createDefaultServer()

	mockServer.Start()
	defer mockServer.Close()

	rtcWSHost := strings.TrimPrefix(mockServer.URL, "http")

	z := createTestRealTimeChatWebsocketService(
		t,
		[]study.RoundTripFunc{
			createSuccessfulChatAuth(t),
		},
		fmt.Sprintf("ws%s", rtcWSHost),
	)

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				testError <- err
				return
			}
		}
	}()

	timeout := time.NewTimer(time.Second * 200)
	select {
	case err := <-testError:
		if err != nil {
			t.Fatal(err)
		}
	case <-timeout.C:
		t.Fatal("did not record a connection within timeout")
		// case successfulConnection := <-mockableZendeskRTCWebsocketServer.state.successfulConnection:
		// 	log.Println("successfulcon check")
		// 	if successfulConnection {
		// 		t.Fatal("expected to fail connection")
		// 	}

		// 	return
	}
}
