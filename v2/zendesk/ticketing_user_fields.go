package zendesk

import "time"

type (
	UserFieldID   uint64
	UserFieldType string
)

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/#json-format
type UserField struct {
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

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/
type TicketingUserFieldsService struct {
	c *client
}
