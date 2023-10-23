package zendesk

import (
	"context"
	"fmt"
	"net/http"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/
type CustomStatusService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#json-format
type CustomStatus struct {
	Active         bool                 `json:"active"`
	ID             CustomStatusID       `json:"id"`
	AgentLabel     string               `json:"agent_label"`
	StatusCategory CustomStatusCategory `json:"status_category"`
}

type CustomStatusCategory string

const (
	New     CustomStatusCategory = "new"
	Open    CustomStatusCategory = "open"
	Pending CustomStatusCategory = "pending"
	Hold    CustomStatusCategory = "hold"
	// Tickets with a "Closed" status belong to the "Solved" status category
	Solved CustomStatusCategory = "solved"
)

type CustomStatusesResponse struct {
	CustomStatuses []CustomStatus `json:"custom_statuses"`
	OffsetPaginationResponse
}

type CustomStatusResponse struct {
	CustomStatus CustomStatus `json:"custom_status"`
}

type CustomStatusPayload struct {
	CustomStatus any `json:"custom_status"`
}

/*
https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#list-custom-ticket-statuses

Does not support cursor pagination
*/
func (s CustomStatusService) List(
	ctx context.Context,
	pageHandler func(response CustomStatusesResponse) error,
) error {
	endpoint := "/api/v2/custom_statuses"

	for {
		target := CustomStatusesResponse{}

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

		if target.NextPage != nil {
			endpoint = *target.NextPage
			continue
		}

		break
	}

	return nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#show-custom-ticket-status
func (s CustomStatusService) Show(ctx context.Context, id CustomStatusID) (CustomStatus, error) {
	target := CustomStatusResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/custom_statuses/%d", id),
		http.NoBody,
	)
	if err != nil {
		return CustomStatus{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return CustomStatus{}, err
	}

	return target.CustomStatus, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/#create-custom-ticket-status
func (s CustomStatusService) Create(ctx context.Context, payload CustomStatusPayload) (CustomStatusResponse, error) {
	target := CustomStatusResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/custom_statuses",
		structToReader(payload),
	)
	if err != nil {
		return CustomStatusResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return CustomStatusResponse{}, err
	}

	return target, nil
}
