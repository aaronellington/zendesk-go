package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UserResponse struct {
	User User `json:"user"`
}

type UserPayload struct {
	User any `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type UsersIncrementalExportResponse struct {
	UsersResponse
	IncrementalExportResponse
}

type User struct {
	ID                   UserID          `json:"id"`
	Active               bool            `json:"active"`
	CreatedAt            time.Time       `json:"created_at"`
	CustomRoleID         *CustomRoleID   `json:"custom_role_id"`
	DefaultGroupID       *GroupID        `json:"default_group_id"`
	Email                string          `json:"email"`
	ExternalID           *string         `json:"external_id"`
	IanaTimeZone         string          `json:"iana_time_zone"`
	LastLoginAt          *time.Time      `json:"last_login_at"`
	Locale               string          `json:"locale"`
	Name                 string          `json:"name"`
	OrganizationID       *OrganizationID `json:"organization_id"`
	Phone                *string         `json:"phone"`
	Role                 string          `json:"role"`
	RoleType             *int            `json:"role_type"`
	Shared               bool            `json:"shared"`
	Suspended            bool            `json:"suspended"`
	Tags                 []string        `json:"tags"`
	TwoFactorAuthEnabled bool            `json:"two_factor_auth_enabled"`
	UpdatedAt            time.Time       `json:"updated_at"`
	Verified             bool            `json:"verified"`
	UserFields           UserFields      `json:"user_fields"`
	Photo                *UserPhoto      `json:"photo"`
}

type UserPhoto struct {
	ContentURL string `json:"content_url"`
}

type UserFields map[string]any

// https://developer.zendesk.com/api-reference/ticketing/users/users/
type UserService struct {
	client *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s UserService) Show(ctx context.Context, id UserID) (User, error) {
	target := UserResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/users/%d", id),
		http.NoBody,
		&target,
	); err != nil {
		return User{}, err
	}

	return target.User, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-self
func (s UserService) ShowSelf(ctx context.Context) (User, error) {
	target := UserResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		"/api/v2/users/me",
		http.NoBody,
		&target,
	); err != nil {
		return User{}, err
	}

	return target.User, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#search-users
func (s UserService) Search(ctx context.Context, query string) (UsersResponse, error) {
	target := UsersResponse{}

	q := url.Values{}
	q.Set("query", query)

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/users/search?%s", q.Encode()),
		http.NoBody,
		&target,
	); err != nil {
		return UsersResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/#incremental-user-export-time-based
func (s UserService) IncrementalExport(
	ctx context.Context,
	startTime int64,
	pageHandler func(response UsersIncrementalExportResponse) error,
) error {
	query := url.Values{}
	query.Set("start_time", fmt.Sprintf("%d", startTime))

	for {
		target := UsersIncrementalExportResponse{}

		if err := s.client.ZendeskRequest(
			ctx,
			http.MethodGet,
			fmt.Sprintf("/api/v2/incremental/users.json?%s", query.Encode()),
			http.NoBody,
			&target,
		); err != nil {
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
func (s UserService) Create(ctx context.Context, payload UserPayload) (UserResponse, error) {
	target := UserResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodPost,
		"/api/v2/users",
		structToReader(payload),
		&target,
	); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#update-user
func (s UserService) Update(ctx context.Context, id UserID, payload UserPayload) (UserResponse, error) {
	target := UserResponse{}

	if err := s.client.ZendeskRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/api/v2/users/%d", id),
		structToReader(payload),
		&target,
	); err != nil {
		return UserResponse{}, err
	}

	return target, nil
}
