package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/
type SatisfactionRatingService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#json-format
type SatisfactionRating struct {
	ID          SatisfactionRatingID `json:"id"`
	AssigneeID  *UserID              `json:"assignee_id"`
	Comment     *string              `json:"comment"`
	CreatedAt   time.Time            `json:"created_at"`
	GroupID     *GroupID             `json:"group_id"`
	Reason      *string              `json:"reason"`
	ReasonCode  *ReasonCode          `json:"reason_code"`
	ReasonID    *ReasonID            `json:"reason_id"`
	RequesterID UserID               `json:"requester_id"`
	Score       string               `json:"score"`
	TicketID    TicketID             `json:"ticket_id"`
	UpdatedAt   *time.Time           `json:"updated_at"`
	URL         string               `json:"url"`
}

type SatisfactionRatingResponse struct {
	SatisfactionRating SatisfactionRating `json:"satisfaction_rating"`
}

type SatisfactionRatingsResponse struct {
	SatisfactionRatings []SatisfactionRating `json:"satisfaction_ratings"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#show-satisfaction-rating
func (s SatisfactionRatingService) Show(ctx context.Context, id SatisfactionRatingID) (SatisfactionRating, error) {
	target := SatisfactionRatingResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/satisfaction_ratings/%d", id),
		http.NoBody,
	)
	if err != nil {
		return SatisfactionRating{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return SatisfactionRating{}, err
	}

	return target.SatisfactionRating, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#list-satisfaction-ratings
func (s SatisfactionRatingService) List(
	ctx context.Context,
	pageHandler func(response SatisfactionRatingsResponse) error,
) error {
	return s.ListWithModifiers(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#filters
func (s SatisfactionRatingService) ListWithModifiers(
	ctx context.Context,
	pageHandler func(response SatisfactionRatingsResponse) error,
	modifiers ...ListTicketSatisfactionRatingModifier,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")

	for _, modifier := range modifiers {
		modifier.ModifyListTicketSatisfactionRatingRequest(&query)
	}

	endpoint := fmt.Sprintf(
		"/api/v2/satisfaction_ratings?%s",
		query.Encode(),
	)

	for {
		target := SatisfactionRatingsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ZendeskRequest(request, &target); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if !target.Meta.HasMore {
			break
		}

		endpoint = target.Links.Next
	}

	return nil
}

type ListTicketSatisfactionRatingModifier interface {
	ModifyListTicketSatisfactionRatingRequest(queryParameters *url.Values)
}

type listTicketSatisfactionRatingModifier func(queryParameters *url.Values)

func (l listTicketSatisfactionRatingModifier) ModifyListTicketSatisfactionRatingRequest(queryParameters *url.Values) {
	l(queryParameters)
}

func WithFilterForStartTime(startTime time.Time) listTicketSatisfactionRatingModifier {
	return listTicketSatisfactionRatingModifier(func(queryParameters *url.Values) {
		queryParameters.Add("start_time", fmt.Sprintf("%d", startTime.Unix()))
	})
}

func WithFilterForEndTime(endTime time.Time) listTicketSatisfactionRatingModifier {
	return listTicketSatisfactionRatingModifier(func(queryParameters *url.Values) {
		queryParameters.Add("end_time", fmt.Sprintf("%d", endTime.Unix()))
	})
}

func WithFilterForScore(score SatisfactionRatingScore) listTicketSatisfactionRatingModifier {
	return listTicketSatisfactionRatingModifier(func(queryParameters *url.Values) {
		queryParameters.Add("score", string(score))
	})
}
