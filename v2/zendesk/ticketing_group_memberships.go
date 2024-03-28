package zendesk

import "time"

type GroupMembershipID uint64

type GroupMembership struct {
	ID        GroupMembershipID `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Default   bool              `json:"default"`
	GroupID   GroupID           `json:"group_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	URL       string            `json:"url"`
	UserID    UserID            `json:"user_id"`
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
type TicketingGroupMembershipsService struct {
	c *client
}
