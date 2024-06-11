package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OrganizationMembershipService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#json-format
type OrganizationMembership struct {
	CreatedAt        time.Time                `json:"created_at"`
	Default          bool                     `json:"default"`
	ID               OrganizationMembershipID `json:"id"`
	OrganizationID   OrganizationID           `json:"organization_id"`
	OrganizationName string                   `json:"organization_name"`
	UpdatedAt        *time.Time               `json:"updated_at"`
	URL              string                   `json:"url"`
	UserID           UserID                   `json:"user_id"`
	ViewTickets      bool                     `json:"view_tickets"`
}

type OrganizationMembershipResponse struct {
	OrganizationMembership OrganizationMembership `json:"organization_membership"`
}

type OrganizationMembershipsResponse struct {
	OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
	CursorPaginationResponse
}

type OrganizationMembershipPayload struct {
	OrganizationMembership OrganizationMembershipPayloadData `json:"organization_membership"`
}

type OrganizationMembershipPayloadData struct {
	UserID         UserID         `json:"user_id"`
	OrganizationID OrganizationID `json:"organization_id"`
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#list-memberships
func (s OrganizationMembershipService) List(
	ctx context.Context,
	pageHandler func(response OrganizationMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/organization_memberships?%s", query.Encode())

	for {
		target := OrganizationMembershipsResponse{}

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

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#list-memberships
func (s OrganizationMembershipService) ListByOrganizationID(
	ctx context.Context,
	organizationID OrganizationID,
	pageHandler func(response OrganizationMembershipsResponse) error,
) error {
	query := url.Values{}
	query.Set("page[size]", "100")
	endpoint := fmt.Sprintf("/api/v2/organizations/%d/organization_memberships?%s", organizationID, query.Encode())

	for {
		target := OrganizationMembershipsResponse{}

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

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#show-membership
func (s OrganizationMembershipService) Show(
	ctx context.Context,
	id OrganizationMembershipID,
) (OrganizationMembership, error) {
	target := OrganizationMembershipResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/organization_memberships/%d", id),
		http.NoBody,
	)
	if err != nil {
		return OrganizationMembership{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return OrganizationMembership{}, err
	}

	return target.OrganizationMembership, nil
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#create-membership
func (s OrganizationMembershipService) Create(
	ctx context.Context,
	payload OrganizationMembershipPayload,
) (OrganizationMembershipResponse, error) {
	target := OrganizationMembershipResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/organization_memberships",
		structToReader(payload),
	)
	if err != nil {
		return OrganizationMembershipResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return OrganizationMembershipResponse{}, err
	}

	return target, nil
}
