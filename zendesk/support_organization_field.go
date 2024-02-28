package zendesk

import (
	"context"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/
type OrganizationFieldService struct {
	client  *client
	generic genericService[
		OrganizationFieldID,
		OrganizationFieldResponse,
		OrganizationFieldsResponse,
	]
}

type OrganizationFieldPayload struct {
	OrganizationField any `json:"organization_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#json-format
type OrganizationField struct {
	Active              bool                  `json:"active"`
	CreatedAt           time.Time             `json:"created_at"`
	CustomFieldOptions  []CustomFieldOption   `json:"custom_field_options"`
	Description         string                `json:"description"`
	ID                  OrganizationFieldID   `json:"id"`
	Key                 string                `json:"key"`
	Position            uint64                `json:"position"`
	RawDescription      string                `json:"raw_description"`
	RawTitle            string                `json:"raw_title"`
	RegexpForValidation *string               `json:"regexp_for_validation"`
	System              bool                  `json:"system"`
	Tag                 Tag                   `json:"tag"`
	Title               string                `json:"title"`
	Type                OrganizationFieldType `json:"type"`
	UpdatedAt           time.Time             `json:"updated_at"`
	URL                 string                `json:"url"`
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
	OrganizationFields []OrganizationField `json:"organization_fields"`
	cursorPaginationResponse
}

type OrganizationFieldResponse struct {
	OrganizationField OrganizationField `json:"organization_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#list-organization-fields
func (s OrganizationFieldService) List(
	ctx context.Context,
	pageHandler func(response OrganizationFieldsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#show-organization-field
func (s OrganizationFieldService) Show(ctx context.Context, id OrganizationFieldID) (OrganizationFieldResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#create-organization-field
func (s OrganizationFieldService) Create(ctx context.Context, payload OrganizationFieldPayload) (OrganizationFieldResponse, error) {
	return s.generic.Create(ctx, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#update-organization-field
func (s OrganizationFieldService) Update(
	ctx context.Context,
	id OrganizationFieldID,
	payload OrganizationFieldPayload,
) (OrganizationFieldResponse, error) {
	return s.generic.Update(ctx, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/#delete-organization-field
func (s OrganizationFieldService) Delete(ctx context.Context, id OrganizationFieldID) error {
	return s.generic.Delete(ctx, id)
}
