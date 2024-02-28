package zendesk

import (
	"context"
	"time"
)

type GroupsResponse struct {
	Groups []Group `json:"groups"`
	cursorPaginationResponse
}

type GroupPayload struct {
	Group any `json:"group"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type Group struct {
	ID        GroupID   `json:"id"`
	IsPublic  bool      `json:"is_public"`
	Name      string    `json:"name"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/
type GroupsService struct {
	client  *client
	generic genericService[
		GroupID,
		GroupResponse,
		GroupsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#list-groups
func (s GroupsService) List(
	ctx context.Context,
	pageHandler func(response GroupsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#show-group
func (s GroupsService) Show(
	ctx context.Context,
	id GroupID,
) (GroupResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/#create-group
func (s GroupsService) Create(ctx context.Context, payload GroupPayload) (GroupResponse, error) {
	return s.generic.Create(ctx, payload)
}
