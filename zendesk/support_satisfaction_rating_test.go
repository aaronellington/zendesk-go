package zendesk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk"
	"github.com/aaronellington/zendesk-go/zendesk/internal/study"
)

func Test_Support_Satisfaction_Ratings_List_200(t *testing.T) {
	ctx := context.Background()

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/satisfaction_rating/list_page1_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/satisfaction_ratings",
				Query: url.Values{
					"page[size]": []string{"100"},
				},
			},
		),
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/satisfaction_rating/list_page2_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   "/api/v2/satisfaction_ratings.json",
				Query: url.Values{
					"page[after]": []string{"aCursor=="},
					"page[size]":  []string{"100"},
				},
			},
		),
	})

	actual := []zendesk.SatisfactionRating{}

	if err := z.Support().SatisfactionRatings().List(
		ctx,
		func(response zendesk.SatisfactionRatingsResponse) error {
			actual = append(actual, response.SatisfactionRatings...)

			return nil
		},
	); err != nil {
		t.Fatal(err)
	}

	if len(actual) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(actual))
	}
}

func Test_Support_Satisfaction_Ratings_Show_200(t *testing.T) {
	ctx := context.Background()

	expectedSatisfactionRatingID := zendesk.SatisfactionRatingID(1000)

	z := createTestService(t, []study.RoundTripFunc{
		study.ServeAndValidate(
			t,
			&study.TestResponseFile{
				StatusCode: http.StatusOK,
				FilePath:   "test_files/responses/support/satisfaction_rating/show_200.json",
			},
			study.ExpectedTestRequest{
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/api/v2/satisfaction_ratings/%d", expectedSatisfactionRatingID),
			},
		),
	})

	assigneeID := zendesk.UserID(9000)
	groupID := zendesk.GroupID(10000)
	comment := "This is a comment!"
	createdTime, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-22T15:28:32Z")

	expectedSatisfactionRating := zendesk.SatisfactionRating{
		URL:         "https://company.zendesk.com/api/v2/satisfaction_ratings/1000.json",
		ID:          expectedSatisfactionRatingID,
		AssigneeID:  &assigneeID,
		GroupID:     &groupID,
		Comment:     &comment,
		CreatedAt:   createdTime,
		RequesterID: zendesk.UserID(9001),
		TicketID:    zendesk.TicketID(99999),
		Score:       string(zendesk.SatisfactionRatingScoreGood),
	}

	actual, err := z.Support().SatisfactionRatings().Show(
		ctx,
		expectedSatisfactionRating.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := study.Assert(expectedSatisfactionRating, actual.SatisfactionRating); err != nil {
		t.Fatal(err)
	}
}
