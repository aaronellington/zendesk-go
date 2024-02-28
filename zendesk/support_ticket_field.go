package zendesk

import (
	"context"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/
type TicketFieldService struct {
	client  *client
	generic genericService[
		TicketFieldID,
		TicketFieldResponse,
		TicketFieldsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#json-format
type TicketField struct {
	Active              bool                `json:"active"`
	AgentDescription    string              `json:"agent_description"`
	CollapsedForAgents  bool                `json:"collapsed_for_agents"`
	CreatedAt           time.Time           `json:"created_at"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	CustomStatuses      []CustomStatus      `json:"custom_statuses"`
	Description         string              `json:"description"`
	EditableInPortal    bool                `json:"editable_in_portal"`
	ID                  TicketFieldID       `json:"id"`
	Position            uint64              `json:"position"`
	RawDescription      string              `json:"raw_description"`
	RawTitle            string              `json:"raw_title"`
	RawTitleInPortal    string              `json:"raw_title_in_portal"`
	RegexpForValidation *string             `json:"regexp_for_validation"`
	Removable           bool                `json:"removable"`
	Required            bool                `json:"required"`
	RequiredInPortal    bool                `json:"required_in_portal"`
	Tag                 Tag                 `json:"tag"`
	Title               string              `json:"title"`
	TitleInPortal       string              `json:"title_in_portal"`
	Type                TicketFieldType     `json:"type"`
	UpdatedAt           time.Time           `json:"updated_at"`
	URL                 string              `json:"url"`
	VisibleInPortal     bool                `json:"visible_in_portal"`
}

type TicketFieldType string

const (
	TicketFieldTypeText              TicketFieldType = "text"              // Default custom field type when type is not specified
	TicketFieldTypeTextArea          TicketFieldType = "textarea"          // For multi-line text
	TicketFieldTypeCheckbox          TicketFieldType = "checkbox"          // To capture a boolean value. Allowed values are true or false
	TicketFieldTypeDate              TicketFieldType = "date"              // Example: 2021-04-16
	TicketFieldTypeInteger           TicketFieldType = "integer"           // String composed of numbers. May contain an optional decimal point
	TicketFieldTypeDecimal           TicketFieldType = "decimal"           // For numbers containing decimals
	TicketFieldTypeRegexp            TicketFieldType = "regexp"            // Matches the Regex pattern found in the custom field settings
	TicketFieldTypePartialCreditCard TicketFieldType = "partialcreditcard" // A credit card number. Only the last 4 digits are retained
	TicketFieldTypeMultiSelect       TicketFieldType = "multiselect"       // Enables users to choose multiple options from a dropdown menu
	TicketFieldTypeTagger            TicketFieldType = "tagger"            // Single-select dropdown menu. It contains one or more tag values belonging to the field's options. Example: ( {"id": 21938362, "value": ["hd_3000", "hd_5555"]})
	TicketFieldTypeLookup            TicketFieldType = "lookup"            // A field to create a relationship  to another object such as a user, ticket, or organization
)

type TicketFieldResponse struct {
	TicketField TicketField `json:"ticket_field"`
}

type TicketFieldsResponse struct {
	TicketFields []TicketField `json:"ticket_fields"`
	cursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#show-ticket-field
func (s TicketFieldService) Show(ctx context.Context, id TicketFieldID) (TicketFieldResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#list-ticket-fields
func (s TicketFieldService) List(
	ctx context.Context,
	pageHandler func(response TicketFieldsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}
