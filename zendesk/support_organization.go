package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OrganizationPayload struct {
	Organization any `json:"organization"`
}

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
	cursorPaginationResponse
}

// NOTE: Organization Fields are returned as a map[string (name of field)]any (value of field), instead of the
// way in which Ticket Fields are returned.
type OrganizationFields map[string]any

func (fields OrganizationFields) GetString(key string) *string {
	rawValue, ok := fields[key]
	if !ok || rawValue == nil {
		return nil
	}

	value, ok := rawValue.(string)
	if !ok {
		panic("organization field " + key + " is not a string")
	}

	return &value
}

func (fields OrganizationFields) GetBool(key string) bool {
	rawValue, ok := fields[key]
	if !ok || rawValue == nil {
		return false
	}

	value, ok := rawValue.(bool)
	if !ok {
		panic("organization field " + key + " is not a string")
	}

	return value
}

type Organization struct {
	ID                 OrganizationID     `json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	DeletedAt          *time.Time         `json:"deleted_at"`
	Details            string             `json:"details"`
	DomainNames        []string           `json:"domain_names"`
	ExternalID         *string            `json:"external_id"`
	GroupID            *GroupID           `json:"group_id"`
	Name               string             `json:"name"`
	Notes              string             `json:"notes"`
	SharedComments     bool               `json:"shared_comments"`
	SharedTickets      bool               `json:"shared_tickets"`
	Tags               []Tag              `json:"tags"`
	UpdatedAt          time.Time          `json:"updated_at"`
	OrganizationFields OrganizationFields `json:"organization_fields"`
}

type OrganizationVia struct {
	Channel string `json:"channel"`
}

type OrganizationSatisfactionRating struct {
	Score string `json:"score"`
}

type OrganizationsIncrementalExportResponse struct {
	OrganizationsResponse
	IncrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
type OrganizationService struct {
	client  *client
	generic genericService[
		OrganizationID,
		OrganizationResponse,
		OrganizationsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#show-organization
func (s OrganizationService) Show(ctx context.Context, id OrganizationID) (OrganizationResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-organization-export
func (s OrganizationService) IncrementalExport(
	ctx context.Context,
	startTime int64,
	pageHandler func(response OrganizationsIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime))

	for {
		target := OrganizationsIncrementalExportResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/organizations.json?%s", query.Encode()),
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

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#create-organization
func (s OrganizationService) Create(ctx context.Context, payload OrganizationPayload) (OrganizationResponse, error) {
	return s.generic.Create(ctx, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#update-organization
func (s OrganizationService) Update(ctx context.Context, id OrganizationID, payload OrganizationPayload) (OrganizationResponse, error) {
	return s.generic.Update(ctx, id, payload)
}

type OrganizationAutocompleteResponse struct {
	Organizations []Organization `json:"organizations"`
	offsetPaginationResponse
}

/*
https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#autocomplete-organizations

Does not support cursor pagination.
*/
func (s OrganizationService) Autocomplete(
	ctx context.Context,
	term string,
	pageHandler func(response OrganizationAutocompleteResponse) error,
) error {
	endpoint := fmt.Sprintf("/api/v2/organizations/autocomplete?name=%s", term)

	for {
		target := OrganizationAutocompleteResponse{}

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
