package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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
}

type UserSearchResponse struct {
	Users           []User            `json:"users"`
	Identities      []UserIdentity    `json:"identities"`
	Organizations   []Organization    `json:"organizations"`
	Groups          []Group           `json:"groups"`
	OpenTicketCount map[string]uint64 `json:"open_ticket_count"`
	OffsetPaginationResponse
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
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s UserService) Show(ctx context.Context, id UserID) (User, error) {
	userInfo, err := s.ShowWithSideloads(ctx, id, nil)
	if err != nil {
		return User{}, err
	}

	return userInfo.User, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s UserService) ShowWithSideloads(
	ctx context.Context,
	id UserID,
	sideloads []UserSideload,
) (UserResponse, error) {
	target := UserResponse{}
	endpoint := fmt.Sprintf("/api/v2/users/%d", id)

	if len(sideloads) > 0 {
		q := url.Values{}

		sideload, sideloads := string(sideloads[0]), sideloads[1:]
		for _, s := range sideloads {
			sideload = fmt.Sprintf("%s,%s", sideload, string(s))
		}

		q.Set("include", sideload)

		endpoint = fmt.Sprintf("%s?%s", endpoint, q.Encode())
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		http.NoBody,
	)
	if err != nil {
		return UserResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-self
func (s UserService) ShowSelf(ctx context.Context) (User, error) {
	target := UserResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/v2/users/me",
		http.NoBody,
	)
	if err != nil {
		return User{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return User{}, err
	}

	return target.User, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#search-users
func (s UserService) Search(ctx context.Context, query string) (UsersResponse, error) {
	target := UsersResponse{}

	q := url.Values{}
	q.Set("query", query)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/users/search?%s", q.Encode()),
		http.NoBody,
	)
	if err != nil {
		return UsersResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UsersResponse{}, err
	}

	return target, nil
}

/*
https://developer.zendesk.com/api-reference/ticketing/users/users/#search-users

Does not support cursor pagination.
*/
func (s UserService) SearchWithSideloads(
	ctx context.Context,
	query string,
	sideloads []UserSideload,
	pageHandler func(response UserSearchResponse) error,
) error {
	q := url.Values{}
	q.Set("query", query)

	if len(sideloads) > 0 {
		sideload, sideloads := string(sideloads[0]), sideloads[1:]
		for _, s := range sideloads {
			sideload = fmt.Sprintf("%s,%s", sideload, string(s))
		}

		q.Set("include", sideload)
	}

	endpoint := fmt.Sprintf("/api/v2/users/search?%s", q.Encode())

	for {
		target := UserSearchResponse{}

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

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s UserService) IncrementalExport(
	ctx context.Context,
	startTime int64,
	pageHandler func(response UsersIncrementalExportResponse) error,
) error {
	return s.IncrementalExportWithSideloads(ctx, startTime, nil, pageHandler)
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s UserService) IncrementalExportWithSideloads(
	ctx context.Context,
	startTime int64,
	sideloads []UserSideload,
	pageHandler func(response UsersIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime))

	if len(sideloads) > 0 {
		sideload, sideloads := string(sideloads[0]), sideloads[1:]
		for _, s := range sideloads {
			sideload = fmt.Sprintf("%s,%s", sideload, string(s))
		}

		query.Set("include", sideload)
	}

	for {
		target := UsersIncrementalExportResponse{}

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/users.json?%s", query.Encode()),
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

		if target.EndOfStream {
			break
		}

		query.Set("start_time", fmt.Sprintf("%d", target.EndTimeUnix))
	}

	return nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#create-user
func (s UserService) Create(
	ctx context.Context,
	payload UserPayload,
) (UserResponse, error) {
	target := UserResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/api/v2/users",
		structToReader(payload),
	)
	if err != nil {
		return UserResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#update-user
func (s UserService) Update(
	ctx context.Context,
	id UserID,
	payload UserPayload,
) (UserResponse, error) {
	target := UserResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/users/%d", id),
		structToReader(payload),
	)
	if err != nil {
		return UserResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}

func (s UserService) DeleteSession(
	ctx context.Context,
	id UserID,
) (UserResponse, error) {
	target := UserResponse{}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/v2/users/%d/sessions", id),
		nil,
	)
	if err != nil {
		return UserResponse{}, err
	}

	if err := s.client.ZendeskRequest(request, &target); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}
