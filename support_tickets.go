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

type TicketID uint64

type TicketResponse struct {
	Ticket Ticket `json:"ticket"`
}

type TicketsResponse struct {
	Tickets []Ticket `json:"tickets"`
}

type Ticket struct {
	AssigneeID         *UserID                  `json:"assignee_id"`
	CreatedAt          time.Time                `json:"created_at"`
	CustomFields       []TicketCustomField      `json:"custom_fields"`
	Description        string                   `json:"description"`
	DueAt              *time.Time               `json:"due_at"`
	ExternalID         *string                  `json:"external_id"`
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

type TicketVia struct {
	Channel string `json:"channel"`
}

type TicketCustomField struct {
	ID    int `json:"id"`
	Value any `json:"value"`
}

type TicketSatisfactionRating struct {
	Score string `json:"score"`
}

type TagsPayload struct {
	Tags Tags `json:"tags"`
}

type Tags []string

type TicketsIncrementalExportResponse struct {
	TicketsResponse
	IncrementalExportResponse
}

type IncrementalExportResponse struct {
	Count       uint64 `json:"count"`
	EndTime     uint64 `json:"end_time"`
	NextPage    string `json:"next_page"`
	EndOfStream bool   `json:"end_of_stream"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
type TicketService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#show-ticket
func (s TicketService) Show(ctx context.Context, id TicketID) (Ticket, error) {
	target := TicketResponse{}

	if err := s.client.jsonRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/tickets/%d", id),
		http.NoBody,
		&target,
	); err != nil {
		return Ticket{}, err
	}

	return target.Ticket, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-ticket-export-time-based
func (s TicketService) IncrementalExport(
	ctx context.Context,
	startTime uint64,
	perPage uint64,
	pageHandler func(response TicketsIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime))
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	for {
		target := TicketsIncrementalExportResponse{}

		if err := s.client.jsonRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/tickets.json?%s", query.Encode()),
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

		query.Set("start_time", fmt.Sprintf("%d", target.EndTime))
	}

	return nil
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

	if err := s.client.jsonRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
		&target,
	); err != nil {
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

	if err := s.client.jsonRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
		&target,
	); err != nil {
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

	if err := s.client.jsonRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/tickets/%d/tags", ticketID),
		payloadBuf,
		&target,
	); err != nil {
		return Tags{}, err
	}

	return target.Tags, nil
}
