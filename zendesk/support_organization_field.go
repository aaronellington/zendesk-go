package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/
type OrganizationFieldService struct {
	client *client
}

type OrganizationFieldPayload struct {
	OrganizationField any `json:"organization_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#json-format
type OrganizationField struct {
	Active      bool                `json:"active"`
	Description string              `json:"description"`
	ID          OrganizationFieldID `json:"id"`
}

type OrganizationFieldsResponse struct {
	OrganizationFields []OrganizationField `json:"organization_fields"`
	CursorPaginationResponse
}

type OrganizationFieldResponse struct {
	OrganizationField OrganizationField `json:"organization_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#list-organization-fields
func (s OrganizationFieldService) List(
	ctx context.Context,
	pageHandler func(response OrganizationFieldsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/organization_fields?%s",
		query.Encode(),
	)

	for {
		target := OrganizationFieldsResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
		)
		if err != nil {
			return err
		}

		if err := s.client.ZendeskRequest(request, &target, true); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#show-organization-field
func (s OrganizationFieldService) Show(ctx context.Context, id OrganizationFieldID) (OrganizationField, error) {
	target := OrganizationFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/organization_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return OrganizationField{}, err
	}

	if err := s.client.ZendeskRequest(request, &target, false); err != nil {
		return OrganizationField{}, err
	}

	return target.OrganizationField, nil
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#create-organization-field
func (s OrganizationFieldService) Create(ctx context.Context, payload OrganizationFieldPayload) (OrganizationFieldResponse, error) {
	target := OrganizationFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/organization_fields",
		structToReader(payload),
	)
	if err != nil {
		return OrganizationFieldResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target, false); err != nil {
		return OrganizationFieldResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#update-organization-field
func (s OrganizationFieldService) Update(
	ctx context.Context,
	id OrganizationFieldID,
	payload OrganizationFieldPayload,
) (OrganizationFieldResponse, error) {
	target := OrganizationFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/organization_fields/%d", id),
		structToReader(payload),
	)
	if err != nil {
		return OrganizationFieldResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target, false); err != nil {
		return OrganizationFieldResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#delete-organization-field
func (s OrganizationFieldService) Delete(ctx context.Context, id OrganizationFieldID) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/organization_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(request, nil, false)
}
