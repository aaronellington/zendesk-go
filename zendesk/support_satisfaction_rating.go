package zendesk

import (
	"context"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/
type SatisfactionRatingService struct {
	client  *client
	generic genericService[
		SatisfactionRatingID,
		SatisfactionRatingResponse,
		SatisfactionRatingsResponse,
	]
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
	cursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#show-satisfaction-rating
func (s SatisfactionRatingService) Show(ctx context.Context, id SatisfactionRatingID) (SatisfactionRatingResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#list-satisfaction-ratings
func (s SatisfactionRatingService) List(
	ctx context.Context,
	pageHandler func(response SatisfactionRatingsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
