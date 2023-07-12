package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type OrganizationID uint64

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
}

type Organization struct {
	ID OrganizationID `json:"id"`
}

type OrganizationVia struct {
	Channel string `json:"channel"`
}

type OrganizationCustomField struct {
	ID    int `json:"id"`
	Value any `json:"value"`
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
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#show-organization
func (s OrganizationService) Show(ctx context.Context, id OrganizationID) (Organization, error) {
	target := OrganizationResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%d", id),
		http.NoBody,
		&target,
	); err != nil {
		return Organization{}, err
	}

	return target.Organization, nil
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

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/organizations.json?%s", query.Encode()),
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
