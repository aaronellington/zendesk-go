package zendesk

// https://developer.zendesk.com/api-reference/ticketing/introduction/
type SupportService struct {
	ticketService *TicketService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func (s *SupportService) Tickets() *TicketService {
	return s.ticketService
}
