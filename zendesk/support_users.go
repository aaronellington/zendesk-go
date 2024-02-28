package zendesk

import (
	"context"
	"time"
)

type UserResponse struct {
	User          User           `json:"user"`
	Identities    []UserIdentity `json:"identities"`
	Organizations []Organization `json:"organizations"`
}

type UserPayload struct {
	User any `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	IncrementalExportResponse

	cursorPaginationResponse
}

type UserSearchResponse struct {
	Users           []User            `json:"users"`
	Identities      []UserIdentity    `json:"identities"`
	Organizations   []Organization    `json:"organizations"`
	Groups          []Group           `json:"groups"`
	OpenTicketCount map[string]uint64 `json:"open_ticket_count"`
	offsetPaginationResponse
}

type UsersIncrementalExportResponse struct {
	UsersResponse
	IncrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#json-format
type User struct {
	ID                   UserID                 `json:"id"`
	Active               bool                   `json:"active"`
	Alias                *string                `json:"alias"`
	OnlyPrivateComments  bool                   `json:"only_private_comments"`
	CreatedAt            time.Time              `json:"created_at"`
	CustomRoleID         *CustomRoleID          `json:"custom_role_id"`
	DefaultGroupID       *GroupID               `json:"default_group_id"`
	Details              *string                `json:"details"`
	Email                string                 `json:"email"`
	ExternalID           *string                `json:"external_id"`
	IanaTimeZone         string                 `json:"iana_time_zone"`
	TimeZone             string                 `json:"time_zone"`
	LastLoginAt          *time.Time             `json:"last_login_at"`
	Locale               string                 `json:"locale"`
	Name                 string                 `json:"name"`
	Notes                *string                `json:"notes"`
	OrganizationID       *OrganizationID        `json:"organization_id"`
	Phone                *string                `json:"phone"`
	RemotePhotoURL       *string                `json:"remote_photo_url"`
	RestrictedAgent      bool                   `json:"restricted_agent"`
	Role                 UserRole               `json:"role"`
	RoleType             *int                   `json:"role_type"`
	Signature            string                 `json:"signature"`
	Shared               bool                   `json:"shared"`
	SharedAgent          bool                   `json:"shared_agent"`
	Suspended            bool                   `json:"suspended"`
	Tags                 []Tag                  `json:"tags"`
	TicketRestriction    *UserTicketRestriction `json:"ticket_restriction"`
	TwoFactorAuthEnabled bool                   `json:"two_factor_auth_enabled"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Verified             bool                   `json:"verified"`
	UserFields           UserFields             `json:"user_fields"`
	Photo                *UserPhoto             `json:"photo"`
}

type UserTicketRestriction string

const (
	UserTicketRestrictionOrganization UserTicketRestriction = "organization"
	UserTicketRestrictionGroups       UserTicketRestriction = "groups"
	UserTicketRestrictionAssigned     UserTicketRestriction = "assigned"
	UserTicketRestrictionRequested    UserTicketRestriction = "requested"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleAgent   UserRole = "agent"
	UserRoleEndUser UserRole = "end-user"
)

type UserPhoto struct {
	ContentURL string `json:"content_url"`
}

// NOTE: User Fields are returned as a map[string (name of field)]any (value of field), instead of the
// way in which Ticket Fields are returned.
type UserFields map[string]any

// https://developer.zendesk.com/api-reference/ticketing/users/users/
type UserService struct {
	client  *client
	generic genericService[
		UserID,
		UserResponse,
		UsersResponse,
	]
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s UserService) Show(ctx context.Context, id UserID) (UserResponse, error) {
	return s.generic.Show(ctx, id)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-self
func (s UserService) ShowSelf(ctx context.Context) (UserResponse, error) {
	return s.generic.getSingle(
		ctx,
		"/api/v2/users/me",
	)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#search-users
func (s UserService) Search(
	ctx context.Context,
	query string,
	pageHandler func(response UsersResponse) error,
) error {
	return s.generic.Search(
		ctx,
		query,
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s UserService) IncrementalExport(
	ctx context.Context,
	startTime int64,
	pageHandler func(response UsersResponse) error,
) error {
	return s.generic.IncrementalExport(
		ctx,
		time.Unix(startTime, 0),
		1000,
		[]string{},
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s UserService) IncrementalExportWithSideloads(
	ctx context.Context,
	startTime int64,
	sideloads []string,
	pageHandler func(response UsersResponse) error,
) error {

	return s.generic.IncrementalExport(
		ctx,
		time.Unix(startTime, 0),
		1000,
		sideloads,
		pageHandler,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#create-user
func (s UserService) Create(
	ctx context.Context,
	payload UserPayload,
) (UserResponse, error) {
	return s.generic.Create(
		ctx,
		payload,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#update-user
func (s UserService) Update(
	ctx context.Context,
	id UserID,
	payload UserPayload,
) (UserResponse, error) {
	return s.generic.Update(
		ctx,
		id,
		payload,
	)
}
