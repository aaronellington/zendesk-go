package zendesk

import (
	"context"
	"time"
)

type ticketingOrganizationObject struct{}

func (r ticketingOrganizationObject) zendeskEntityName() string {
	return "organizations"
}

type OrganizationID int64

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#json-format
type Organization struct {
	ID                 OrganizationID          `json:"id"`
	CreatedAt          time.Time               `json:"created_at"`
	DeletedAt          *time.Time              `json:"deleted_at"`
	Details            string                  `json:"details"`
	DomainNames        []string                `json:"domain_names"`
	ExternalID         *string                 `json:"external_id"`
	GroupID            *GroupID                `json:"group_id"`
	Name               string                  `json:"name"`
	Notes              string                  `json:"notes"`
	SharedComments     bool                    `json:"shared_comments"`
	SharedTickets      bool                    `json:"shared_tickets"`
	Tags               []Tag                   `json:"tags"`
	UpdatedAt          time.Time               `json:"updated_at"`
	OrganizationFields OrganizationFieldValues `json:"organization_fields"`
}

type OrganizationFieldValues map[string]any

func (fields OrganizationFieldValues) GetString(key string) string {
	rawValue, ok := fields[key]
	if !ok || rawValue == nil {
		return ""
	}

	value, ok := rawValue.(string)
	if !ok {
		return ""
	}

	return value
}

func (fields OrganizationFieldValues) GetBool(key string) bool {
	rawValue, ok := fields[key]
	if !ok || rawValue == nil {
		return false
	}

	value, ok := rawValue.(bool)
	if !ok {
		return false
	}

	return value
}

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
	ticketingOrganizationObject
}

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
	ticketingOrganizationObject
	cursorPaginationResponse
}

type OrganizationPayload struct {
	Organization any `json:"organization"`
}

type OrganizationsIncrementalExportResponse struct {
	Organizations []Organization `json:"organizations"`
	ticketingOrganizationObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
type TicketingOrganizationsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#create-organization
func (s *TicketingOrganizationsService) Create(
	ctx context.Context,
	payload OrganizationPayload,
) (OrganizationResponse, error) {
	return createRequest[OrganizationResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#show-organization
func (s *TicketingOrganizationsService) Show(
	ctx context.Context,
	id OrganizationID,
) (OrganizationResponse, error) {
	return showRequest[OrganizationID, OrganizationResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#list-organizations
func (s *TicketingOrganizationsService) List(
	ctx context.Context,
	pageHandler func(response OrganizationsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#update-organization
func (s *TicketingOrganizationsService) Update(
	ctx context.Context,
	id OrganizationID,
	payload OrganizationPayload,
) (OrganizationResponse, error) {
	return updateRequest[OrganizationID, OrganizationResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/#delete-organization
func (s *TicketingOrganizationsService) Delete(
	ctx context.Context,
	id OrganizationID,
) error {
	return deleteRequest[OrganizationID, OrganizationResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-organization-export
func (s *TicketingOrganizationsService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(OrganizationsIncrementalExportResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return incrementalExportRequest(ctx, s.c, startTime, pageHandler, requestQueryModifiers...)
}
