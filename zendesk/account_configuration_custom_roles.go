package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles
type CustomRoleService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#json-format
type CustomRole struct {
	Configuration   CustomRoleConfiguration `json:"configuration"`
	CreatedAt       time.Time               `json:"created_at"`
	Description     string                  `json:"description"`
	ID              CustomRoleID            `json:"id"`
	Name            string                  `json:"name"`
	RoleType        uint64                  `json:"role_type"`
	TeamMemberCount uint64                  `json:"team_member_count"`
	UdpatedAt       time.Time               `json:"updated_at"`
}

// A list of custom object keys mapped to JSON objects that define the agent's permissions (scopes) for each object.
// Allowed values: "read", "update", "delete", "create". The "read" permission is required if any other scopes are
// specified. Example: { "shipment": { "scopes": ["read", "update"] } }
type CustomRoleCustomObject struct {
	Scopes []AgentScopePermission `json:"scopes"`
}

type AgentScopePermission string

const (
	ScopeRead   AgentScopePermission = "read"
	ScopeUpdate AgentScopePermission = "update"
	ScopeDelete AgentScopePermission = "delete"
	ScopeCreate AgentScopePermission = "create"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#configuration
type CustomRoleConfiguration struct {
	AssignTicketsToAnyGroup bool                              `json:"assign_tickets_to_any_group"`
	ChatAccess              bool                              `json:"chat_access"`
	CustomObjects           map[string]CustomRoleCustomObject `json:"custom_objects"`
}

type CustomRolesResponse struct {
	CustomRoles []CustomRole `json:"custom_roles"`
}

type CustomRoleResponse struct {
	CustomRole CustomRole `json:"custom_role"`
}

/*
https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#list-custom-roles

Does not support pagination
*/
func (s CustomRoleService) List(
	ctx context.Context,
) ([]CustomRole, error) {
	target := CustomRolesResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/custom_roles",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return nil, err
	}

	return target.CustomRoles, nil
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#show-custom-role
func (s CustomRoleService) Show(ctx context.Context, id CustomRoleID) (CustomRole, error) {
	target := CustomRoleResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/custom_roles/%d", id),
		http.NoBody,
	)
	if err != nil {
		return CustomRole{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return CustomRole{}, err
	}

	return target.CustomRole, nil
}
