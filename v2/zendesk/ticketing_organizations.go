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

type Organization struct {
	ID OrganizationID `json:"id"`
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
