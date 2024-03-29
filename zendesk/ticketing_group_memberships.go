package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type GroupMembershipID uint64

type GroupMembership struct {
	ID        GroupMembershipID `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Default   bool              `json:"default"`
	GroupID   GroupID           `json:"group_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	URL       string            `json:"url"`
	UserID    UserID            `json:"user_id"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
type TicketingGroupMembershipsService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#create-membership
func (s *TicketingGroupMembershipsService) Create(
	ctx context.Context,
	payload GroupMembershipPayload,
) (GroupResponse, error) {
	return createRequest[GroupResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#show-membership
func (s *TicketingGroupMembershipsService) Show(
	ctx context.Context,
	id GroupMembershipID,
) (GroupMembershipResponse, error) {
	return showRequest[GroupMembershipID, GroupMembershipResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s *TicketingGroupMembershipsService) List(
	ctx context.Context,
	pageHandler func(response GroupMembershipsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s *TicketingGroupMembershipsService) ListByUser(
	ctx context.Context,
	userID UserID,
	pageHandler func(response GroupMembershipsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return paginatedRequest(
		s.c,
		ctx,
		fmt.Sprintf("/api/v2/users/%d/group_memberships", userID),
		pageHandler,
		requestQueryModifiers...,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s *TicketingGroupMembershipsService) ListByGroup(
	ctx context.Context,
	groupID GroupID,
	pageHandler func(response GroupMembershipsResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return paginatedRequest(
		s.c,
		ctx,
		fmt.Sprintf("/api/v2/groups/%d/memberships", groupID),
		pageHandler,
		requestQueryModifiers...,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#delete-membership
func (s *TicketingGroupMembershipsService) Delete(
	ctx context.Context,
	id GroupMembershipID,
) error {
	return deleteRequest[GroupMembershipID, GroupMembershipResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#set-membership-as-default
func (s *TicketingGroupMembershipsService) SetDefault(
	ctx context.Context,
	userID UserID,
	groupMembershipID GroupMembershipID,
) (GroupMembershipsResponse, error) {
	return genericRequest[GroupMembershipsResponse](
		s.c,
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/users/%d/group_memberships/%d/make_default", userID, groupMembershipID),
		http.NoBody,
	)
}

type GroupMembershipResponse struct {
	GroupMembership GroupMembership `json:"group_membership"`
	ticketingGroupMembershipObject
}

type GroupMembershipsResponse struct {
	GroupMemberships []GroupMembership `json:"group_memberships"`
	ticketingGroupMembershipObject
	cursorPaginationResponse
}

type GroupMembershipPayload struct {
	GroupMembership any `json:"group_membership"`
}

type ticketingGroupMembershipObject struct{}

func (r ticketingGroupMembershipObject) zendeskEntityName() string {
	return "group_memberships"
}
