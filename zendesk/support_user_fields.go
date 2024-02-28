package zendesk

import (
	"context"
	"time"
)

type UserFieldResponse struct {
	UserField UserFieldConfiguration `json:"user_field"`
}

type UserFieldsResponse struct {
	UserFields []UserFieldConfiguration `json:"user_fields"`
	cursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#json-format
type UserFieldConfiguration struct {
	Active              bool                `json:"active"`
	CreatedAt           time.Time           `json:"created_at"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	Description         *string             `json:"description"`
	ID                  UserFieldID         `json:"id"`
	Key                 string              `json:"key"`
	Position            uint64              `json:"position"`
	RawDescription      *string             `json:"raw_description"`
	RawTitle            *string             `json:"raw_title"`
	RegexpForValidation *string             `json:"regexp_for_validation"`
	System              bool                `json:"system"`
	Tag                 Tag                 `json:"tag"`
	Title               *string             `json:"title"`
	Type                UserFieldType       `json:"type"`
	UpdatedAt           *time.Time          `json:"updated_at"`
	URL                 string              `json:"url"`
}

type UserFieldType string

const (
	UserFieldTypeText     UserFieldType = "text"     // Default custom field type when type is not specified
	UserFieldTypeTextArea UserFieldType = "textarea" // For multi-line text
	UserFieldTypeCheckbox UserFieldType = "checkbox" // To capture a boolean value. Allowed values are true or false
	UserFieldTypeDate     UserFieldType = "date"     // Example: 2021-04-16
	UserFieldTypeDropdown UserFieldType = "dropdown" //
	UserFieldTypeInteger  UserFieldType = "integer"  // String composed of numbers. May contain an optional decimal point
	UserFieldTypeDecimal  UserFieldType = "decimal"  // For numbers containing decimals
	UserFieldTypeRegexp   UserFieldType = "regexp"   // Matches the Regex pattern found in the custom field settings
	UserFieldTypeLookup   UserFieldType = "lookup"   // A field to create a relationship  to another object such as a user, ticket, or organization
)

type CustomFieldOption struct {
	ID       CustomFieldOptionID `json:"id"`
	Name     string              `json:"name"`
	Position uint64              `json:"position"`
	RawName  string              `json:"raw_name"`
	URL      string              `json:"url"`
	Value    string              `json:"value"`
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/
type UserFieldService struct {
	client  *client
	generic genericService[
		UserFieldID,
		UserFieldResponse,
		UserFieldsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#list-user-fields
func (s UserFieldService) List(
	ctx context.Context,
	pageHandler func(response UserFieldsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#show-user-field
func (s UserFieldService) Show(
	ctx context.Context,
	id UserFieldID,
) (UserFieldResponse, error) {
	return s.generic.Show(ctx, id)
}
