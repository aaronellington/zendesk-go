package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/
type OrganizationFieldService struct {
	client *client
}

type OrganizationFieldPayload struct {
	OrganizationField any `json:"organization_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#json-format
type OrganizationFieldConfiguration struct {
	Active              bool                `json:"active"`
	CreatedAt           time.Time           `json:"created_at"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	Description         string              `json:"description"`
	ID                  OrganizationFieldID `json:"id"`
	Key                 string              `json:"key"`
	Position            uint64              `json:"position"`
	RawDescription      string              `json:"raw_description"`
	RawTitle            string              `json:"raw_title"`
	RegexpForValidation *string             `json:"regexp_for_validation"`
	// RelationshipFilter object `json:"relationship_filter"`
	// RelationshipTargetType string `json:"relationship_target_type"`
	System    bool                  `json:"system"`
	Tag       Tag                   `json:"tag"`
	Title     string                `json:"title"`
	Type      OrganizationFieldType `json:"type"`
	UpdatedAt time.Time             `json:"updated_at"`
	URL       string                `json:"url"`
}

type OrganizationFieldType string

const (
	OrganizationFieldTypeText     OrganizationFieldType = "text"     // Default custom field type when type is not specified
	OrganizationFieldTypeTextArea OrganizationFieldType = "textarea" // For multi-line text
	OrganizationFieldTypeCheckbox OrganizationFieldType = "checkbox" // To capture a boolean value. Allowed values are true or false
	OrganizationFieldTypeDate     OrganizationFieldType = "date"     // Example: 2021-04-16
	OrganizationFieldTypeDropdown OrganizationFieldType = "dropdown" //
	OrganizationFieldTypeInteger  OrganizationFieldType = "integer"  // String composed of numbers. May contain an optional decimal point
	OrganizationFieldTypeDecimal  OrganizationFieldType = "decimal"  // For numbers containing decimals
	OrganizationFieldTypeRegexp   OrganizationFieldType = "regexp"   // Matches the Regex pattern found in the custom field settings
	OrganizationFieldTypeLookup   OrganizationFieldType = "lookup"   // A field to create a relationship  to another object such as a user, ticket, or organization
)

type OrganizationFieldsResponse struct {
	OrganizationFields []OrganizationFieldConfiguration `json:"organization_fields"`
	CursorPaginationResponse
}

type OrganizationFieldResponse struct {
	OrganizationField OrganizationFieldConfiguration `json:"organization_field"`
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

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#show-organization-field
func (s OrganizationFieldService) Show(ctx context.Context, id OrganizationFieldID) (OrganizationFieldConfiguration, error) {
	target := OrganizationFieldResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/organization_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return OrganizationFieldConfiguration{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return OrganizationFieldConfiguration{}, err
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

	if err := s.client.ZendeskRequest(request, &target); err != nil {
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

	if err := s.client.ZendeskRequest(request, &target); err != nil {
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

	return s.client.ZendeskRequest(request, nil)
}
