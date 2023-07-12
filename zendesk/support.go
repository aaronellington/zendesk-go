package zendesk

// https://developer.zendesk.com/api-reference/ticketing/introduction/
type SupportService struct {
	organizationService *OrganizationService
	userService         *UserService
	ticketService       *TicketService
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *SupportService) Organizations() *OrganizationService {
	return s.organizationService
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/
func (s *SupportService) Users() *UserService {
	return s.userService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func (s *SupportService) Tickets() *TicketService {
	return s.ticketService
}
