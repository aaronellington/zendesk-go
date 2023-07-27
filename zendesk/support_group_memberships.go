package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GroupMembershipPayload struct {
	GroupMembership any `json:""`
}

type GroupMembership struct {
	ID        GroupMembershipID `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Default   bool              `json:"default"`
	GroupID   GroupID           `json:"group_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	UserID    UserID            `json:"user_id"`
}

type GroupMembershipsResponse struct {
	GroupMemberships []GroupMembership
	CursorPaginationResponse
}

type GroupMembershipResponse struct {
	GroupMembership GroupMembership `json:"group_membership"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
type GroupMembershipService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) List(
	ctx context.Context,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/group_memberships?%s", query.Encode())

	for {
		target := GroupMembershipsResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
			&target,
		); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) ListByUser(
	ctx context.Context,
	userID UserID,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/users/%d/group_memberships?%s", userID, query.Encode())

	for {
		target := GroupMembershipsResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
			&target,
		); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#list-memberships
func (s GroupMembershipService) ListByGroup(
	ctx context.Context,
	groupID GroupID,
	pageHandler func(response GroupMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/groups/%d/memberships?%s", groupID, query.Encode())

	for {
		target := GroupMembershipsResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			endpoint,
			http.NoBody,
			&target,
		); err != nil {
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

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#set-membership-as-default
func (s GroupMembershipService) SetDefault(
	ctx context.Context,
	userID UserID,
	groupMembershipID GroupMembershipID,
) (GroupMembershipsResponse, error) {
	target := GroupMembershipsResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/users/%d/group_memberships/%d/make_default", userID, groupMembershipID),
		http.NoBody,
		&target,
	); err != nil {
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
	target := GroupMembershipResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodPost,
		"/api/v2/group_memberships",
		structToReader(GroupMembershipPayload{
			GroupMembership: map[string]any{
				"userID":  userID,
				"groupID": groupID,
			},
		}),
		nil,
	); err != nil {
		return GroupMembershipResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#delete-membership
func (s GroupMembershipService) Delete(
	ctx context.Context,
	userID UserID,
	groupMembershipID GroupMembershipID,
) error {
	return s.client.ZendeskRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/users/%d/group_memberships/%d", userID, groupMembershipID),
		http.NoBody,
		nil,
	)
}
