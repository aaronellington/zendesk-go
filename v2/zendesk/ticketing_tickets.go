package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ticketsTicketObject struct{}

func (r ticketsTicketObject) zendeskEntityName() string {
	return "tickets"
}

type TicketID uint64

type Ticket struct {
	ID TicketID `json:"id"`
}

type TicketResponse struct {
	Ticket Ticket `json:"ticket"`
	ticketsTicketObject
}

type TicketsResponse struct {
	Tickets []Ticket `json:"tickets"`
	ticketsTicketObject
	cursorPaginationResponse
}

type TicketPayload struct {
	Ticket any `json:"ticket"`
}

type TicketsIncrementalExportResponse struct {
	Tickets []Ticket `json:"tickets"`
	ticketsTicketObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
type TicketingTicketsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#create-ticket
func (s *TicketingTicketsService) Create(
	ctx context.Context,
	payload TicketPayload,
) (TicketResponse, error) {
	return createRequest[TicketResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#show-ticket
func (s *TicketingTicketsService) Show(
	ctx context.Context,
	id TicketID,
) (TicketResponse, error) {
	return showRequest[TicketID, TicketResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#show-ticket
func (s *TicketingTicketsService) List(
	ctx context.Context,
	pageHandler func(response TicketsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#update-ticket
func (s *TicketingTicketsService) Update(
	ctx context.Context,
	id TicketID,
	payload TicketPayload,
) (TicketResponse, error) {
	return updateRequest[TicketID, TicketResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#update-ticket
func (s *TicketingTicketsService) Delete(
	ctx context.Context,
	id TicketID,
) error {
	return deleteRequest[TicketID, TicketResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-ticket-export-time-based
func (s *TicketingTicketsService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(TicketsIncrementalExportResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return incrementalExportRequest(ctx, s.c, startTime, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#list-ticket-incidents
func (s *TicketingTicketsService) ListIncidents(
	ctx context.Context,
	id TicketID,
	pageHandler func(response TicketsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return paginatedRequest(
		s.c,
		ctx,
		fmt.Sprintf("/api/v2/tickets/%d/incidents", id),
		pageHandler,
		requestQueryModifiers...,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#add-tags
func (s *TicketingTicketsService) AddTags(
	ctx context.Context,
	id TicketID,
	tags []Tag,
) (TagsResponse, error) {
	return genericRequest[TagsResponse](
		s.c,
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/tickets/%d/tags", id),
		TagsPayload{
			Tags: tags,
		},
	)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (s *TicketingTicketsService) RemoveTags(
	ctx context.Context,
	id TicketID,
	tags []Tag,
) (TagsResponse, error) {
	return genericRequest[TagsResponse](
		s.c,
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/tickets/%d/tags", id),
		TagsPayload{
			Tags: tags,
		},
	)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (s *TicketingTicketsService) SetTags(
	ctx context.Context,
	id TicketID,
	tags []Tag,
) (TagsResponse, error) {
	return genericRequest[TagsResponse](
		s.c,
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/tags", id),
		TagsPayload{
			Tags: tags,
		},
	)
}
