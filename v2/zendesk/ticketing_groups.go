package zendesk

import (
	"context"
)

type ticketingGroupObject struct{}

func (r ticketingGroupObject) zendeskEntityName() string {
	return "groups"
}

type GroupID int64

type Group struct {
	ID GroupID `json:"id"`
}

type GroupResponse struct {
	Group Group `json:"group"`
	ticketingGroupObject
}

type GroupsResponse struct {
	Groups []Group `json:"groups"`
	ticketingGroupObject
	cursorPaginationResponse
}

type GroupPayload struct {
	Group any `json:"group"`
}

type GroupsIncrementalExportResponse struct {
	Groups []Group `json:"groups"`
	ticketingGroupObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/
type TicketingGroupsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#create-group
func (s *TicketingGroupsService) Create(
	ctx context.Context,
	payload GroupPayload,
) (GroupResponse, error) {
	return createRequest[GroupResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#show-group
func (s *TicketingGroupsService) Show(
	ctx context.Context,
	id GroupID,
) (GroupResponse, error) {
	return showRequest[GroupID, GroupResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#list-groups
func (s *TicketingGroupsService) List(
	ctx context.Context,
	pageHandler func(response GroupsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#update-group
func (s *TicketingGroupsService) Update(
	ctx context.Context,
	id GroupID,
	payload GroupPayload,
) (GroupResponse, error) {
	return updateRequest[GroupID, GroupResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#delete-group
func (s *TicketingGroupsService) Delete(
	ctx context.Context,
	id GroupID,
) error {
	return deleteRequest[GroupID, GroupResponse](ctx, s.c, id)
}
