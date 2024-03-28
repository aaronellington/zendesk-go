package zendesk

import (
	"context"
	"math"
	"time"
)

type ticketingScheduleObject struct{}

func (r ticketingScheduleObject) zendeskEntityName() string {
	return "business_hours/schedules"
}

type ScheduleID uint64

type Schedule struct {
	ID        ScheduleID         `json:"id"`
	Name      string             `json:"name"`
	TimeZone  string             `json:"time_zone"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Intervals []ScheduleInterval `json:"intervals"`
}

func (schedule Schedule) Location() (*time.Location, error) {
	timeZoneLabel := schedule.TimeZone

	switch schedule.TimeZone {
	case "Eastern Time (US & Canada)":
		timeZoneLabel = "America/New_York"
	case "Central Time (US & Canada)":
		timeZoneLabel = "America/Chicago"
	case "Pacific Time (US & Canada)":
		timeZoneLabel = "America/Los_Angeles"
	}

	loc, err := time.LoadLocation(timeZoneLabel)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func (schedule Schedule) Active(now time.Time) (bool, error) {
	loc, err := schedule.Location()
	if err != nil {
		return false, err
	}

	localTime := now.In(loc)
	sunday := localTime.Add(time.Hour * -24 * time.Duration(localTime.Weekday()))
	beginningOfWeek := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, loc)
	minutesIntoTheWeek := int(math.Floor(now.Sub(beginningOfWeek).Minutes()))

	for _, interval := range schedule.Intervals {
		if minutesIntoTheWeek >= interval.StartTime && minutesIntoTheWeek < interval.EndTime {
			return true, nil
		}
	}

	return false, nil
}

type ScheduleInterval struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
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
