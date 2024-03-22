package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func TestRealTimeChatRest_GetAllChatMetrics_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_all_chat_metrics_200.json",
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

	expectedChatDuration := int64(325)

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

func TestRealTimeChatRest_GetAllChatMetricsForDepartment_200(t *testing.T) {
	ctx := context.Background()

	departmentID := zendesk.GroupID(9000)

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_all_chat_metrics_for_department_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/stream/chats",
				Query: url.Values{
					"department_id": []string{strconv.FormatUint(uint64(departmentID), 10)},
				},
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetAllChatMetricsForDepartment(ctx, departmentID)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				MissedChats: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
					SixtyMinuteWindow:  0,
				},
				ActiveChats:   0,
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
			},
			DepartmentID: &departmentID,
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}

func TestRealTimeChatRest_GetAllChatMetricsForDepartment_WithNegativeValue_200(t *testing.T) {
	ctx := context.Background()

	departmentID := zendesk.GroupID(9000)

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat/get_all_chat_metrics_for_department_200_With_Negative_Response.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/stream/chats",
				Query: url.Values{
					"department_id": []string{strconv.FormatUint(uint64(departmentID), 10)},
				},
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetAllChatMetricsForDepartment(ctx, departmentID)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				MissedChats: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
					SixtyMinuteWindow:  0,
				},
				ActiveChats:   -1,
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
			},
			DepartmentID: &departmentID,
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}

func TestRealTimeChatRest_GetAllChatMetricsForSpecificTimeWindow_200(t *testing.T) {
	ctx := context.Background()

	timeWindow := zendesk.LiveChatTimeWindow30Minutes

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_all_chat_metrics_for_time_window_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/stream/chats",
				Query: url.Values{
					"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
				},
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetAllChatMetricsForSpecificTimeWindow(ctx, timeWindow)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				MissedChats: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
				},
				ActiveChats:   0,
				IncomingChats: 0,
				AssignedChats: 0,
				SatisfactionBad: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
				},
				SatisfactionGood: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
				},
			},
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}

func TestRealTimeChatRest_GetSingleChatMetric_200(t *testing.T) {
	ctx := context.Background()

	liveChatMetricKey := zendesk.LiveChatMetricKeyIncomingChats

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_single_chat_metric_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/stream/chats/%s", liveChatMetricKey),
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetSingleChatMetric(ctx, liveChatMetricKey)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				IncomingChats: 0,
			},
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}

func TestRealTimeChatRest_GetSingleChatMetricForDepartment_200(t *testing.T) {
	ctx := context.Background()

	liveChatMetricKey := zendesk.LiveChatMetricKeyIncomingChats
	departmentID := zendesk.GroupID(9000)

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_single_chat_metric_for_department_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/stream/chats/%s", liveChatMetricKey),
				Query: url.Values{
					"department_id": []string{strconv.FormatUint(uint64(departmentID), 10)},
				},
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetSingleChatMetricForDepartment(ctx, liveChatMetricKey, departmentID)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				IncomingChats: 0,
			},
			DepartmentID: &departmentID,
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}

func TestRealTimeChatRest_GetSingleChatMetricForSpecificTimeWindow_200(t *testing.T) {
	ctx := context.Background()

	liveChatMetricKey := zendesk.LiveChatMetricKeyMissedChats
	timeWindow := zendesk.LiveChatTimeWindow30Minutes

	z := createTestService(t, []study.RoundTripFunc{
		createSuccessfulChatAuth(t),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/livechat/realtimechat_rest/get_single_chat_metric_for_time_window_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/stream/chats/%s", liveChatMetricKey),
				Query: url.Values{
					"window": []string{strconv.FormatUint(uint64(timeWindow), 10)},
				},
			},
		),
	})

	actual, err := z.LiveChat().RealTimeChat().RealTimeChatRestService().GetSingleChatMetricForSpecificTimeWindow(ctx, liveChatMetricKey, timeWindow)
	if err != nil {
		t.Fatal(err)
	}

	expected := zendesk.ChatsStreamResponse{
		StatusCode: 200,
		Content: zendesk.ChatsStreamResponseContent{
			Topic: "chats",
			Type:  "update",
			Data: zendesk.ChatMetrics{
				MissedChats: &zendesk.ChatMetricWindow{
					ThirtyMinuteWindow: 0,
				},
			},
		},
	}

	if err := study.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}
