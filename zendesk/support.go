package zendesk

// https://developer.zendesk.com/api-reference/ticketing/introduction/
type SupportService struct {
	customStatusService           *CustomStatusService
	groupMembershipService        *GroupMembershipService
	groupService                  *GroupsService
	organizationFieldService      *OrganizationFieldService
	organizationMembershipService *OrganizationMembershipService
	organizationService           *OrganizationService
	satisfactionRatingService     *SatisfactionRatingService
	scheduleService               *ScheduleService
	suspendedTicketService        *SuspendedTicketService
	ticketAttachmentService       *TicketAttachmentService
	ticketAuditService            *TicketAuditService
	ticketCommentService          *TicketCommentService
	ticketFieldService            *TicketFieldService
	ticketFormService             *TicketFormService
	ticketService                 *TicketService
	ticketTagService              *TicketTagService
	userFieldsService             *UserFieldService
	userService                   *UserService
	userIdentityService           *UserIdentityService
	sideConversationService       *SideConversationService
	automationService             *AutomationService
	triggerService                *TriggerService
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *SupportService) CustomStatuses() *CustomStatusService {
	return s.customStatusService
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *SupportService) OrganizationFields() *OrganizationFieldService {
	return s.organizationFieldService
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *SupportService) OrganizationMemberships() *OrganizationMembershipService {
	return s.organizationMembershipService
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

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings
func (s *SupportService) SatisfactionRatings() *SatisfactionRatingService {
	return s.satisfactionRatingService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
func (s *SupportService) TicketForms() *TicketFormService {
	return s.ticketFormService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func (s *SupportService) Tickets() *TicketService {
	return s.ticketService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_audits/
func (s *SupportService) TicketAudits() *TicketAuditService {
	return s.ticketAuditService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/
func (s *SupportService) TicketAttachments() *TicketAttachmentService {
	return s.ticketAttachmentService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments
func (s *SupportService) TicketComments() *TicketCommentService {
	return s.ticketCommentService
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields
func (s *SupportService) TicketFields() *TicketFieldService {
	return s.ticketFieldService
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/
func (s *SupportService) AutomationService() *AutomationService {
	return s.automationService
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/
func (s *SupportService) TriggerService() *TriggerService {
	return s.triggerService
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags
func (s *SupportService) TicketTags() *TicketTagService {
	return s.ticketTagService
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
