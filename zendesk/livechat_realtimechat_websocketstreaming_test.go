package zendesk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestRealTimeChatWebsocketStreaming_Connect_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat/get_all_chat_metrics_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/stream/chats",
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetAllChatMetrics(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedChatDuration := uint64(325)

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				MissedChats: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 1,
					SixtyMinuteWindow:  1,
				},
				ActiveChats:   1,
				IncomingChats: 0,
				AssignedChats: 0,
				SatisfactionBad: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
					SixtyMinuteWindow:  0,
				},
				SatisfactionGood: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
					SixtyMinuteWindow:  0,
				},
				ChatDurationAvg: &expectedChatDuration,
				ChatDurationMax: &expectedChatDuration,
			},
			DepartmentID: nil,
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}
