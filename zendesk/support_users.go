package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	UserID uint64
)

func (userID *UserID) UnmarshalJSON(b []byte) error {
	// Try it as a uint64 first
	var targetUint64 uint64
	if err := json.Unmarshal(b, &targetUint64); err == nil {
		*userID = UserID(targetUint64)

		return nil
	}

	// Only try it as a string as a last resort
	var targetString string
	if err := json.Unmarshal(b, &targetString); err != nil {
		return err
	}

	typeUint64, _ := strconv.ParseUint(targetString, 0, 64)
	*userID = UserID(typeUint64)

	return nil
}

type UserResponse struct {
	User User `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type UsersIncrementalExportResponse struct {
	UsersResponse
	IncrementalExportResponse
}

type User struct {
	ID UserID `json:"id"`
}

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

		query.Set("start_time", fmt.Sprintf("%d", target.EndTime))
	}

	return nil
}
