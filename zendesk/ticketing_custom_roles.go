package zendesk

type CustomRoleID uint64

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
type TicketingCustomRolesService struct {
	c *client
}
