package zendesk

import (
	"context"
	"math"
	"net/http"
	"time"
)

type SchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
	ID        ScheduleID         `json:"id"`
	Name      string             `json:"name"`
	TimeZone  string             `json:"time_zone"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Intervals []ScheduleInterval `json:"intervals"`
}

type ScheduleInterval struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

func (schedule Schedule) Location() (*time.Location, error) {
	timeZoneLabel := schedule.TimeZone

	switch schedule.TimeZone {
	case "Eastern Time (US & Canada)":
		timeZoneLabel = "America/New_York"
	case "Central Time (US & Canada)":
		timeZoneLabel = "America/Chicago"
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

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/
type ScheduleService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/#list-schedules
func (s *ScheduleService) List(
	ctx context.Context,
) (SchedulesResponse, error) {
	target := SchedulesResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		"/api/v2/business_hours/schedules",
		http.NoBody,
		&target,
	); err != nil {
		return SchedulesResponse{}, err
	}

	return target, nil
}
