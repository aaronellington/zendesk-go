package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GroupMembershipPayload struct {
	GroupMembership any `json:"group_membership"`
}

type GroupMembership struct {
	ID        GroupMembershipID `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Default   bool              `json:"default"`
	GroupID   GroupID           `json:"group_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	URL       string            `json:"url"`
	UserID    UserID            `json:"user_id"`
}

type GroupMembershipsResponse struct {
	GroupMemberships []GroupMembership `json:"group_memberships"`
	cursorPaginationResponse
}

type GroupMembershipResponse struct {
	GroupMembership GroupMembership `json:"group_membership"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
type GroupMembershipService struct {
	client  *client
	generic genericService[
		GroupMembershipID,
		GroupMembershipResponse,
		GroupMembershipsResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) List(
	ctx context.Context,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) ListByUser(
	ctx context.Context,
	userID UserID,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/users/%d/group_memberships?%s", userID, query.Encode())

	return genericList(
		ctx,
		s.client,
		endpoint,
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) ListByGroup(
	ctx context.Context,
	groupID GroupID,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/groups/%d/memberships?%s", groupID, query.Encode())

	return genericList(
		ctx,
		s.client,
		endpoint,
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#set-membership-as-default
func (s GroupMembershipService) SetDefault(
	ctx context.Context,
	userID UserID,
	groupMembershipID GroupMembershipID,
) (GroupMembershipsResponse, error) {
	target := GroupMembershipsResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/users/%d/group_memberships/%d/make_default", userID, groupMembershipID),
		http.NoBody,
	)
	if err != nil {
		return GroupMembershipsResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return GroupMembershipsResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#create-membership
func (s GroupMembershipService) Create(
	ctx context.Context,
	userID UserID,
	groupID GroupID,
) (GroupMembershipResponse, error) {
	return s.generic.Create(
		ctx,
		GroupMembershipPayload{
			GroupMembership: map[string]any{
				"user_id":  userID,
				"group_id": groupID,
			},
		},
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#delete-membership
func (s GroupMembershipService) Delete(
	ctx context.Context,
	id GroupMembershipID,
) error {
	return s.generic.Delete(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#show-membership
func (s GroupMembershipService) Show(
	ctx context.Context,
	id GroupMembershipID,
) (GroupMembershipResponse, error) {
	return s.generic.Show(ctx, id)
}
