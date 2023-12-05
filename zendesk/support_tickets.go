package zendesk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TicketPayload struct {
	Ticket any `json:"ticket"`
}

type TicketResponse struct {
	Ticket Ticket `json:"ticket"`
}

type TicketsResponse struct {
	Tickets []Ticket `json:"tickets"`
}

type MergeRequestPayload struct {
	IDs                   []TicketID `json:"ids"`
	SourceComment         string     `json:"source_comment"`
	SourceCommentIsPublic bool       `json:"source_comment_is_public"`
	TargetComment         string     `json:"target_comment"`
	TargetCommentIsPublic bool       `json:"target_comment_is_public"`
}

type Ticket struct {
	AssigneeID         *UserID                  `json:"assignee_id"`
	CreatedAt          time.Time                `json:"created_at"`
	CollaboratorIDs    []UserID                 `json:"collaborator_ids"`
	CustomFields       []TicketField            `json:"custom_fields"`
	Description        string                   `json:"description"`
	DueAt              *time.Time               `json:"due_at"`
	ExternalID         *string                  `json:"external_id"`
	Fields             []TicketField            `json:"fields"`
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

	Dates TicketDates `json:"dates"`
}

type TicketDates struct {
	AssigneeUpdatedAt    *time.Time `json:"assignee_updated_at"`
	RequesterUpdatedAt   *time.Time `json:"requester_updated_at"`
	StatusUpdatedAt      *time.Time `json:"status_updated_at"`
	InitiallyAssignedAt  *time.Time `json:"initially_updated_at"`
	AssignedAt           *time.Time `json:"assigned_updated_at"`
	SolvedAt             *time.Time `json:"solved_updated_at"`
	LatestCommentAddedAt *time.Time `json:"latest_comment_added_at"`
}

type TicketVia struct {
	Channel string `json:"channel"`
}

type TicketFields []TicketField

func (fields TicketFields) CreateMap() map[TicketFieldID]any {
	fieldMap := map[TicketFieldID]any{}
	for _, field := range fields {
		fieldMap[field.ID] = field.Value
	}

	return fieldMap
}

type TicketField struct {
	ID    TicketFieldID `json:"id"`
	Value any           `json:"value"`
}

type TicketSatisfactionRating struct {
	Score string `json:"score"`
}

type TagsPayload struct {
	Tags Tags `json:"tags"`
}

type Tags []Tag

func (tags Tags) HasTag(targetTag Tag) bool {
	for _, tag := range tags {
		if tag == targetTag {
			return true
		}
	}

	return false
}

type TicketsIncrementalExportResponse struct {
	TicketsResponse
	IncrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
type TicketService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#show-ticket
func (s TicketService) Show(ctx context.Context, id TicketID) (Ticket, error) {
	target := TicketResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/tickets/%d", id),
		http.NoBody,
	)
	if err != nil {
		return Ticket{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Ticket{}, err
	}

	return target.Ticket, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#create-ticket
func (s TicketService) Create(ctx context.Context, payload TicketPayload) (TicketResponse, error) {
	target := TicketResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/tickets",
		structToReader(payload),
	)
	if err != nil {
		return TicketResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_import/#ticket-import
func (s TicketService) Import(ctx context.Context, payload TicketPayload) (TicketResponse, error) {
	target := TicketResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/imports/tickets",
		structToReader(payload),
	)
	if err != nil {
		return TicketResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#merge-tickets-into-target-ticket
func (s TicketService) Merge(ctx context.Context, destination TicketID, payload MergeRequestPayload) (JobStatusResponse, error) {
	target := JobStatusResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/merge", destination),
		structToReader(payload),
	)
	if err != nil {
		return JobStatusResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return JobStatusResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#update-ticket
func (s TicketService) Update(ctx context.Context, id TicketID, payload TicketPayload) (TicketResponse, error) {
	target := TicketResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/tickets/%d", id),
		structToReader(payload),
	)
	if err != nil {
		return TicketResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-ticket-export-time-based
func (s TicketService) IncrementalExportWithSideloads(
	ctx context.Context,
	startTime time.Time,
	sideloads []TicketSideload,
	perPage uint,
	pageHandler func(response TicketsIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime.Unix()))
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	if len(sideloads) > 0 {
		sideload, sideloads := string(sideloads[0]), sideloads[1:]
		for _, s := range sideloads {
			sideload = fmt.Sprintf("%s,%s", sideload, string(s))
		}

		query.Set("include", sideload)
	}

	for {
		target := TicketsIncrementalExportResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/tickets.json?%s", query.Encode()),
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

		if target.EndOfStream {
			break
		}

		query.Set("start_time", fmt.Sprintf("%d", target.EndTimeUnix))
	}

	return nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-ticket-export-time-based
func (s TicketService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	perPage uint,
	pageHandler func(response TicketsIncrementalExportResponse) error,
) error {
	return s.IncrementalExportWithSideloads(ctx, startTime, nil, perPage, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#add-tags
func (s TicketService) AddTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#set-tags
func (s TicketService) SetTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#remove-tags
func (s TicketService) RemoveTags(ctx context.Context, ticketID TicketID, tags Tags) (Tags, error) {
	target := TagsPayload{}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(TagsPayload{
		Tags: tags,
	}); err != nil {
		return Tags{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
	)
	if err != nil {
		return Tags{}, err
	}

	if err := s.client.ZendeskRequest(
		request,
		&target,
	); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}

type ListProblemTicketIncidentsResponse struct {
	Tickets []Ticket `json:"tickets"`
	CursorPaginationResponse
}

func (s TicketService) ListProblemTicketIncidents(
	ctx context.Context,
	problemTicket TicketID,
	pageHandler func(response ListProblemTicketIncidentsResponse) error,
) error {
	query := url.Values{}
	// Default values
	query.Set("page[size]", "100")

	endpoint := fmt.Sprintf("/api/v2/tickets/%d/incidents.json?%s", problemTicket, query.Encode())

	for {
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		target := ListProblemTicketIncidentsResponse{}

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
