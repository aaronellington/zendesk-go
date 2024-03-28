package zendesk

import (
	"context"
	"net/http"
	"time"
)

type ticketingUserObject struct{}

func (r ticketingUserObject) zendeskEntityName() string {
	return "users"
}

type (
	UserID   int64
	UserRole string
)

type UserTicketRestriction string

const (
	UserTicketRestrictionOrganization UserTicketRestriction = "organization"
	UserTicketRestrictionGroups       UserTicketRestriction = "groups"
	UserTicketRestrictionAssigned     UserTicketRestriction = "assigned"
	UserTicketRestrictionRequested    UserTicketRestriction = "requested"
)

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
	UserFields           UserFieldValues        `json:"user_fields"`
	Photo                *UserPhoto             `json:"photo"`
}

type UserPhoto struct {
	ContentURL string `json:"content_url"`
}

type UserFieldValues map[string]any

func (fields UserFieldValues) GetString(key string) string {
	rawValue, ok := fields[key]
	if !ok || rawValue == nil {
		return ""
	}

	value, ok := rawValue.(string)
	if !ok {
		return ""
	}

	return value
}

type UserResponse struct {
	User User `json:"user"`
	ticketingUserObject
}

type UsersResponse struct {
	Users []User `json:"users"`
	ticketingUserObject
	cursorPaginationResponse
}

type UserPayload struct {
	User any `json:"user"`
}

type UsersIncrementalExportResponse struct {
	Users []User `json:"users"`
	ticketingUserObject
	incrementalExportResponse
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/
type TicketingUsersService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#create-user
func (s *TicketingUsersService) Create(
	ctx context.Context,
	payload UserPayload,
) (UserResponse, error) {
	return createRequest[UserResponse](ctx, s.c, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s *TicketingUsersService) Show(
	ctx context.Context,
	id UserID,
) (UserResponse, error) {
	return showRequest[UserID, UserResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-self
func (s *TicketingUsersService) ShowSelf(
	ctx context.Context,
) (UserResponse, error) {
	return genericRequest[UserResponse](
		s.c,
		ctx,
		http.MethodGet,
		"/api/v2/users/me",
		http.NoBody,
	)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#list-users
func (s *TicketingUsersService) List(
	ctx context.Context,
	pageHandler func(response UsersResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return listRequest(ctx, s.c, pageHandler, requestQueryModifiers...)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#update-user
func (s *TicketingUsersService) Update(
	ctx context.Context,
	id UserID,
	payload UserPayload,
) (UserResponse, error) {
	return updateRequest[UserID, UserResponse](ctx, s.c, id, payload)
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#delete-user
func (s *TicketingUsersService) Delete(
	ctx context.Context,
	id UserID,
) error {
	return deleteRequest[UserID, UserResponse](ctx, s.c, id)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s *TicketingUsersService) IncrementalExport(
	ctx context.Context,
	startTime time.Time,
	pageHandler func(UsersIncrementalExportResponse) error,
	requestQueryModifiers ...RequestQueryModifiers,
) error {
	return incrementalExportRequest(ctx, s.c, startTime, pageHandler, requestQueryModifiers...)
}
