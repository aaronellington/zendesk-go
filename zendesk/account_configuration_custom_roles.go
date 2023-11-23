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
	UpdatedAt       time.Time               `json:"updated_at"`
}

// A list of custom object keys mapped to JSON objects that define the agent's permissions (scopes) for each object.
// Allowed values: "read", "update", "delete", "create". The "read" permission is required if any other scopes are
// specified. Example: { "shipment": { "scopes": ["read", "update"] } }.
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
	AssignTicketsToAnyGroup      bool                              `json:"assign_tickets_to_any_group"`
	ChatAccess                   bool                              `json:"chat_access"`
	CustomObjects                map[string]CustomRoleCustomObject `json:"custom_objects"`
	EndUserListAccess            string                            `json:"end_user_list_access"`
	EndUserProfileAccess         string                            `json:"end_user_profile_access"`
	ExploreAccess                string                            `json:"explore_access"`
	ForumAccess                  string                            `json:"forum_access"`
	ForumAccessRestrictedContent bool                              `json:"forum_access_restricted_content"`
	GroupAccess                  bool                              `json:"group_access"`
	LightAgent                   bool                              `json:"light_agent"`
	MacroAccess                  string                            `json:"macro_access"`
	ManageAutomations            bool                              `json:"manage_automations"`
	ManageBusinessRules          bool                              `json:"manage_business_rules"`
	ManageContextualWorkspaces   bool                              `json:"manage_contextual_workspaces"`
	ManageDynamicContent         bool                              `json:"manage_dynamic_content"`
	ManageExtensionsAndChannels  bool                              `json:"manage_extensions_and_channels"`
	ManageFacebook               bool                              `json:"manage_facebook"`
	ManageGroupMemberships       bool                              `json:"manage_group_memberships"`
	ManageGroups                 bool                              `json:"manage_groups"`
	ManageOrganizationFields     bool                              `json:"manage_organization_fields"`
	ManageOrganizations          bool                              `json:"manage_organizations"`
	ManageRoles                  string                            `json:"manage_roles"`
	ManageSkills                 bool                              `json:"manage_skills"`
	ManageSLAs                   bool                              `json:"manage_slas"`
	ManageSuspendedTickets       bool                              `json:"manage_suspended_tickets"`
	ManageTeamMembers            string                            `json:"manage_team_members"`
	ManageTicketFields           bool                              `json:"manage_ticket_fields"`
	ManageTicketForms            bool                              `json:"manage_ticket_forms"`
	ManageTriggers               bool                              `json:"manage_triggers"`
	ManageUserFields             bool                              `json:"manage_user_fields"`
	ModerateForums               bool                              `json:"moderate_forums"`
	OrganizationEditing          bool                              `json:"organization_editing"`
	OrganizationNotesEditing     bool                              `json:"organization_notes_editing"`
	ReportAccess                 string                            `json:"report_access"`
	SideConversationCreate       bool                              `json:"side_conversation_create"`
	TicketAccess                 string                            `json:"ticket_access"`
	TicketCommentAccess          string                            `json:"ticket_comment_access"`
	TicketDeletion               bool                              `json:"ticket_deletion"`
	TicketEditing                bool                              `json:"ticket_editing"`
	TicketMerge                  bool                              `json:"ticket_merge"`
	TicketRedaction              bool                              `json:"ticket_redaction"`
	TicketTagEditing             bool                              `json:"ticket_tag_editing"`
	TwitterSearchAccess          bool                              `json:"twitter_search_access"`
	UserViewAccess               string                            `json:"user_view_access"`
	ViewAccess                   string                            `json:"view_access"`
	ViewDeletedTickets           bool                              `json:"view_deleted_tickets"`
	VoiceAccess                  bool                              `json:"voice_access"`
	VoiceDashboardAccess         bool                              `json:"voice_dashboard_access"`
}

type CustomRolesResponse struct {
	CustomRoles []CustomRole `json:"custom_roles"`
	OffsetPaginationResponse
}

type CustomRoleResponse struct {
	CustomRole CustomRole `json:"custom_role"`
}

/*
https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#list-custom-roles

Does not support cursor pagination.
*/
func (s CustomRoleService) List(
	ctx context.Context,
	pageHandler func(response CustomRolesResponse) error,
) error {
	endpoint := "/api/v2/custom_roles"

	for {
		target := CustomRolesResponse{}

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

		if target.NextPage != nil {
			endpoint = *target.NextPage

			continue
		}

		break
	}

	return nil
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
