package zendesk

import "time"

type (
	TicketFieldID   uint64
	TicketFieldType string
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#json-format
type TicketField struct {
	ID                  TicketFieldID       `json:"id"`
	Active              bool                `json:"active"`
	AgentDescription    string              `json:"agent_description"`
	CollapsedForAgents  bool                `json:"collapsed_for_agents"`
	CreatedAt           time.Time           `json:"created_at"`
	CustomFieldOptions  []CustomFieldOption `json:"custom_field_options"`
	Description         string              `json:"description"`
	EditableInPortal    bool                `json:"editable_in_portal"`
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

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/
type TicketingTicketFieldsService struct {
	c *client
}
