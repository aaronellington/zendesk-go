package zendesk

import "context"

type TicketActivityID uint64

// https://developer.zendesk.com/api-reference/ticketing/tickets/activity_stream/#json-format
type TicketActivity struct {
	ID TicketActivityID `json:"id"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/activity_stream/
type TicketingActivitiesService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/activity_stream/#show-activity
func (s *TicketingActivitiesService) Show(
	ctx context.Context,
	id TicketActivityID,
	requestQueryModifiers ...RequestQueryModifiers,
) (TicketActivityResponse, error) {
	return showRequest[TicketActivityID, TicketActivityResponse](
		ctx,
		s.c,
		id,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/activity_stream/#list-activities
func (s *TicketingActivitiesService) List(
	ctx context.Context,
	pageHandler func(response TicketActivitiesResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(
		ctx,
		s.c,
		pageHandler,
	)
}

type TicketActivityResponse struct {
	Activity TicketActivity `json:"activity"`
	ticketingTicketActivityObject
}

type TicketActivitiesResponse struct {
	Activities []TicketActivity `json:"activities"`
	ticketingTicketActivityObject
	cursorPaginationResponse
}

type ticketingTicketActivityObject struct{}

func (r ticketingTicketActivityObject) zendeskEntityName() string {
	return "activities"
}
