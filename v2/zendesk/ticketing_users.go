package zendesk

import (
	"context"
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

type User struct {
	ID UserID `json:"id"`
}

type UserPhoto struct {
	ContentURL string `json:"content_url"`
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
