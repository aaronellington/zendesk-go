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
	target := GroupMembershipResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/group_memberships",
		structToReader(GroupMembershipPayload{
			GroupMembership: map[string]any{
				"user_id":  userID,
				"group_id": groupID,
			},
		}),
	)
	if err != nil {
		return GroupMembershipResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
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
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/users/%d/group_memberships/%d", userID, groupMembershipID),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	return s.client.ZendeskRequest(
		request,
		nil,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/#show-membership
func (s GroupMembershipService) Show(
	ctx context.Context,
	id GroupMembershipID,
) (GroupMembership, error) {
	target := GroupMembershipResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/group_memberships/%d", id),
		http.NoBody,
	)
	if err != nil {
		return GroupMembership{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return GroupMembership{}, err
	}

	return target.GroupMembership, nil
}
