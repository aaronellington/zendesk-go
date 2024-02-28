package zendesk

import (
	"context"
	"time"
)

type OrganizationMembershipService struct {
	client  *client
	generic genericService[
		OrganizationMembershipID,
		OrganizationMembershipResponse,
		OrganizationMembershipsResponse,
	]
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
	cursorPaginationResponse
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
	return s.generic.List(ctx, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#show-membership
func (s OrganizationMembershipService) Show(
	ctx context.Context,
	id OrganizationMembershipID,
) (OrganizationMembershipResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#create-membership
func (s OrganizationMembershipService) Create(
	ctx context.Context,
	payload OrganizationMembershipPayload,
) (OrganizationMembershipResponse, error) {
	return s.generic.Create(ctx, payload)
}
