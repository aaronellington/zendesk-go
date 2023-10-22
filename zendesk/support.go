package zendesk

// https://developer.zendesk.com/api-reference/ticketing/introduction/
type SupportService struct {
	groupMembershipService  *GroupMembershipService
	groupService            *GroupsService
	organizationService     *OrganizationService
	scheduleService         *ScheduleService
	suspendedTicketService  *SuspendedTicketService
	ticketAuditService      *TicketAuditService
	ticketService           *TicketService
	userFieldsService       *UserFieldService
	userService             *UserService
	userIdentityService     *UserIdentityService
	sideConversationService *SideConversationService
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *SupportService) Organizations() *OrganizationService {
	return s.organizationService
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/
func (s *SupportService) Groups() *GroupsService {
	return s.groupService
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
func (s *SupportService) GroupMemberships() *GroupMembershipService {
	return s.groupMembershipService
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/
func (s *SupportService) Users() *UserService {
	return s.userService
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/
func (s *SupportService) UserIdentities() *UserIdentityService {
	return s.userIdentityService
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/
func (s *SupportService) UserFields() *UserFieldService {
	return s.userFieldsService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func (s *SupportService) Tickets() *TicketService {
	return s.ticketService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_audits/
func (s *SupportService) TicketAudits() *TicketAuditService {
	return s.ticketAuditService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/
func (s *SupportService) SuspendedTickets() *SuspendedTicketService {
	return s.suspendedTicketService
}

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation/
func (s *SupportService) SideConversations() *SideConversationService {
	return s.sideConversationService
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/
func (s *SupportService) Schedules() *ScheduleService {
	return s.scheduleService
}
