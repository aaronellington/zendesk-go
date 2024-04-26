package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/
type TicketFieldService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#json-format
type TicketFieldConfiguration struct {
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

type TicketFieldConfigurationResponse struct {
	TicketField TicketFieldConfiguration `json:"ticket_field"`
}

type TicketFieldsConfigurationResponse struct {
	TicketFields []TicketFieldConfiguration `json:"ticket_fields"`
	CursorPaginationResponse
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#show-ticket-field
func (s TicketFieldService) Show(ctx context.Context, id TicketFieldID) (TicketFieldConfiguration, error) {
	target := TicketFieldConfigurationResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/ticket_fields/%d", id),
		http.NoBody,
	)
	if err != nil {
		return TicketFieldConfiguration{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketFieldConfiguration{}, err
	}

	return target.TicketField, nil
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#list-ticket-fields
func (s TicketFieldService) List(
	ctx context.Context,
	pageHandler func(response TicketFieldsConfigurationResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf(
		"/api/v2/ticket_fields?%s",
		query.Encode(),
	)

	for {
		target := TicketFieldsConfigurationResponse{}

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

type TicketFieldPayload struct {
	TicketField any `json:"ticket_field"`
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#update-ticket-field
func (s TicketFieldService) Update(ctx context.Context, id TicketFieldID, payload TicketFieldPayload) (TicketFieldConfigurationResponse, error) {
	target := TicketFieldConfigurationResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/ticket_fields/%d", id),
		structToReader(payload),
	)
	if err != nil {
		return TicketFieldConfigurationResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return TicketFieldConfigurationResponse{}, err
	}

	return target, nil
}
