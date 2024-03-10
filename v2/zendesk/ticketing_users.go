package zendesk

import "context"

type UserID uint64

type User struct {
	ID UserID `json:"id"`
}

type UserResponse struct {
	User User `json:"user"`
}

func (r UserResponse) zendeskEntityName() string {
	return "users"
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/
type TicketingUsersService struct {
	c *client
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (s *TicketingUsersService) Show(ctx context.Context, id UserID) (UserResponse, error) {
	return showRequest[UserID, UserResponse](ctx, s.c, id)
}
