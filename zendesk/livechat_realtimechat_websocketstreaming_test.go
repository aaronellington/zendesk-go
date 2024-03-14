package zendesk_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestTest(t *testing.T) {

	ctx := context.Background()

	z, mockRealTimeChatWebsocketServer := createTestRealTimeChatWebsocketService(
		t,
		ctx,
		settings{},
	)

	success := make(chan error)
	testCaseChan := make(chan error)

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAgentMetric(ctx, zendesk.LiveChatMetricKeyAgentsOnline); err != nil {
			success <- err
			return
		}

		ticker := time.NewTicker(time.Millisecond * 250)
		for {
			select {
			case <-testCaseChan:
				return
			case <-ticker.C:
				_ = mockRealTimeChatWebsocketServer.history

				metrics, err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetAgentMetric(ctx)
				if err != nil {
					// TODO: log
					log.Println(err)
					continue
					// not ready
				}

				log.Println(metrics)
				// TODO: get

				// validate data we got back
				// if notValid {
				// TODO: log

				// 	continue
				// }
				ticker.Stop()
				success <- errors.New("This is a fake error!")
				return
			}
		}

	}()

	timeout := time.NewTimer(time.Second * 5)
	select {
	case err := <-mockRealTimeChatWebsocketServer.connError:
		testCaseChan <- err

		t.Fatal(err)
	case <-timeout.C:
		t.Fatal("took too long")
	case err := <-success:
		log.Printf("%+v", mockRealTimeChatWebsocketServer.history)
		t.Fatal(err)
		return
	}

	_ = mockRealTimeChatWebsocketServer.conn.Close()
	// Make no goroutes are running
}

func TestTest2(t *testing.T) {
	ctx := context.Background()

	z, mockRealTimeChatWebsocketServer := createTestRealTimeChatWebsocketService(
		t,
		ctx,
		settings{},
	)

	success := make(chan error)

	go func() {
		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().SubscribeToAgentMetric(ctx, zendesk.LiveChatMetricKeyAgentsOnline); err != nil {
			success <- err
			return
		}

		for range time.NewTicker(time.Second).C {
			_ = mockRealTimeChatWebsocketServer.history

			metrics, err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().GetAgentMetric(ctx)
			if err != nil {
				// TODO: log
				log.Println(err)
				continue
				// not ready
			}

			log.Println(metrics)
			// TODO: get

			// validate data we got back
			// if notValid {
			// TODO: log

			// 	continue
			// }

			success <- errors.New("This is a fake error!")
			break
		}
	}()

	timeout := time.NewTimer(time.Second * 2)
	select {
	case err := <-mockRealTimeChatWebsocketServer.connError:
		t.Fatal(err)
	case <-timeout.C:
		t.Fatal("took too long")
	case err := <-success:
		log.Printf("%+v", mockRealTimeChatWebsocketServer.history)
		t.Fatal(err)
		return
	}

	_ = mockRealTimeChatWebsocketServer.conn.Close()
	// Make no goroutes are running
}

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	testError := make(chan error)

	z, mockRealTimeChatWebsocketServer := createTestRealTimeChatWebsocketService(
		t,
		[]study.RoundTripFunc{
			createSuccessfulChatAuth(t),
		},
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
	for {
		select {
		case err := <-testError:
			if err != nil {
				t.Fatal(err)
			}

		case <-timeout.C:
			t.Fatal("did not record a connection within timeout")
		}
	}
}

// func TestRealTimeChatWebsocketStreaming_Connect_401(t *testing.T) {
// 	ctx := context.Background()

// 	testError := make(chan error)

// 	mockRealTimeChatWebsocketServer := newMockRealTimeChatWebsocketServer(
// 		testError,
// 		nil,
// 		&settings{
// 			ValidOAuthToken: "No Valid Token",
// 		},
// 		frameHandlers{},
// 	)

// 	mockServer := mockRealTimeChatWebsocketServer.createDefaultServer()

// 	mockServer.Start()
// 	defer mockServer.Close()

// 	rtcWSHost := strings.Replace(mockServer.URL, "http", "ws", 1)

// 	z := createTestRealTimeChatWebsocketService(
// 		t,
// 		[]study.RoundTripFunc{
// 			createSuccessfulChatAuth(t),
// 		},
// 		rtcWSHost,
// 	)

// 	go func() {
// 		if err := z.LiveChat().RealTimeChat().RealTimeChatStreamingService().ConnectToWebsocket(ctx); err != nil {
// 			if !errors.Is(err, context.Canceled) {
// 				testError <- err
// 				return
// 			}
// 		}
// 	}()

// 	timeout := time.NewTimer(time.Second * 2)
// 	select {
// 	case err := <-testError:
// 		if err == nil {
// 			t.Fatal("expected to receive an error")
// 		}
// 		for _, message := range mockRealTimeChatWebsocketServer.history.sentFrames {
// 			fmt.Printf("%s\n", string(message.payload))
// 		}
// 	case <-timeout.C:
// 		t.Fatal("did not record a connection within timeout")
// 	}
// }

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
