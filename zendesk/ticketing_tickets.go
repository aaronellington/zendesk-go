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

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#json-format
type Ticket struct {
	AssigneeID         *UserID                  `json:"assignee_id"`
	CreatedAt          time.Time                `json:"created_at"`
	CollaboratorIDs    []UserID                 `json:"collaborator_ids"`
	CustomFields       TicketFieldValues        `json:"custom_fields"`
	Description        string                   `json:"description"`
	DueAt              *time.Time               `json:"due_at"`
	ExternalID         *string                  `json:"external_id"`
	Fields             TicketFieldValues        `json:"fields"`
	FollowerIDs        []UserID                 `json:"follower_ids"`
	GroupID            *GroupID                 `json:"group_id"`
	HasIncidents       bool                     `json:"has_incidents"`
	ID                 TicketID                 `json:"id"`
	IsPublic           bool                     `json:"is_public"`
	OrganizationID     *OrganizationID          `json:"organization_id"`
	Priority           string                   `json:"priority"`
	ProblemID          *TicketID                `json:"problem_id"`
	RequesterID        UserID                   `json:"requester_id"`
	SatisfactionRating TicketSatisfactionRating `json:"satisfaction_rating"`
	Status             string                   `json:"status"`
	Subject            string                   `json:"subject"`
	SubmitterID        UserID                   `json:"submitter_id"`
	Tags               Tags                     `json:"tags"`
	TicketFormID       TicketFormID             `json:"ticket_form_id"`
	Type               *string                  `json:"type"`
	UpdatedAt          time.Time                `json:"updated_at"`
	URL                string                   `json:"url"`
	Via                TicketVia                `json:"via"`
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

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#list-tickets
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

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#delete-ticket
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

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#merge-tickets-into-target-ticket
func (s *TicketingTicketsService) Merge(
	ctx context.Context,
	destination TicketID,
	payload MergeRequestPayload,
) (JobStatusResponse, error) {
	return genericRequest[JobStatusResponse](
		s.c,
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/merge", destination),
		payload,
	)
}

type TicketFieldValue struct {
	ID    TicketFieldID `json:"id"`
	Value any           `json:"value"`
}

type TicketFieldValues []TicketFieldValue

func (fields TicketFieldValues) CreateMap() map[TicketFieldID]any {
	fieldMap := map[TicketFieldID]any{}
	for _, field := range fields {
		fieldMap[field.ID] = field.Value
	}

	return fieldMap
}

type TicketSatisfactionRating struct {
	Score string `json:"score"`
}

type TicketVia struct {
	Channel string `json:"channel"`
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

type MergeRequestPayload struct {
	IDs                   []TicketID `json:"ids"`
	SourceComment         string     `json:"source_comment"`
	SourceCommentIsPublic bool       `json:"source_comment_is_public"`
	TargetComment         string     `json:"target_comment"`
	TargetCommentIsPublic bool       `json:"target_comment_is_public"`
}
