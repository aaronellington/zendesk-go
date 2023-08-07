package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TicketMetricEventsResponse struct {
	TicketMetricEvents []TicketMetricEvent `json:"ticket_metric_events"`
}

type TicketMetricEventsIncrementalExportResponse struct {
	TicketMetricEventsResponse
	IncrementalExportResponse
}

type TicketMetricEvent struct {
	ID         TicketMetricEventID    `json:"id"`
	TicketID   TicketID               `json:"ticket_id"`
	InstanceID uint64                 `json:"instance_id"`
	Metric     string                 `json:"metric"`
	Type       string                 `json:"type"`
	Time       time.Time              `json:"time"`
	SLA        *ServiceLevelAgreement `json:"sla,omitempty"`
}

type ServiceLevelAgreement struct {
	Target        uint64    `json:"target"`
	BusinessHours bool      `json:"business_hours"`
	Policy        SLAPolicy `json:"policy"`
}

type SLAPolicy struct {
	ID          SLAPolicyID `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_metric_events
type TicketMetricEventService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_metric_events/#list-ticket-metric-events
func (s TicketMetricEventService) List(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(response TicketMetricEventsIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime.Unix()))
	query.Set("page[size]", "100")

	for {
		target := TicketMetricEventsIncrementalExportResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/ticket_metric_events?%s", query.Encode()),
			http.NoBody,
			&target,
		); err != nil {
			return err
		}

		if err := pageHandler(target); err != nil {
			return err
		}

		if target.EndOfStream {
			break
		}

		query.Set("start_time", fmt.Sprintf("%d", target.EndTimeUnix))
	}

	return nil
}
