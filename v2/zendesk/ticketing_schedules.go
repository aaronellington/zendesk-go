package zendesk

import (
	"context"
)

type ticketingScheduleObject struct{}

func (r ticketingScheduleObject) zendeskEntityName() string {
	return "business_hours/schedules"
}

type ScheduleID int64

type Schedule struct {
	ID ScheduleID `json:"id"`
}

type ScheduleResponse struct {
	Schedule Schedule `json:"schedule"`
	ticketingScheduleObject
}

type SchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
	ticketingScheduleObject
	cursorPaginationResponse
}

type SchedulePayload struct {
	Schedule any `json:"schedule"`
}

type SchedulesIncrementalExportResponse struct {
	Schedules []Schedule `json:"schedules"`
	ticketingScheduleObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/
type TicketingSchedulesService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#create-schedule
func (s *TicketingSchedulesService) Create(
	ctx context.Context,
	payload SchedulePayload,
) (ScheduleResponse, error) {
	return createRequest[ScheduleResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#show-schedule
func (s *TicketingSchedulesService) Show(
	ctx context.Context,
	id ScheduleID,
) (ScheduleResponse, error) {
	return showRequest[ScheduleID, ScheduleResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#list-schedules
func (s *TicketingSchedulesService) List(
	ctx context.Context,
	pageHandler func(response SchedulesResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#update-schedule
func (s *TicketingSchedulesService) Update(
	ctx context.Context,
	id ScheduleID,
	payload SchedulePayload,
) (ScheduleResponse, error) {
	return updateRequest[ScheduleID, ScheduleResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#delete-schedule
func (s *TicketingSchedulesService) Delete(
	ctx context.Context,
	id ScheduleID,
) error {
	return deleteRequest[ScheduleID, ScheduleResponse](ctx, s.c, id)
}
