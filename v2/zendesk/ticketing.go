package zendesk

type TicketingService struct {
	accountSettings              *TicketingAccountSettingsService
	activityStream               *TicketingActivityStreamService
	appLocationInstallations     *TicketingAppLocationInstallationsService
	appLocations                 *TicketingAppLocationsService
	apps                         *TicketingAppsService
	auditLogs                    *TicketingAuditLogsService
	automations                  *TicketingAutomationsService
	bookmarks                    *TicketingBookmarksService
	brands                       *TicketingBrandsService
	customRoles                  *TicketingCustomRolesService
	customTicketStatuses         *TicketingCustomTicketStatusesService
	groupMemberships             *TicketingGroupMembershipsService
	groups                       *TicketingGroupsService
	groupSLAPolicies             *TicketingGroupSLAPoliciesService
	incrementalExports           *TicketingIncrementalExportsService
	incrementalSkillBasedRouting *TicketingIncrementalSkillBasedRoutingService
	jobStatuses                  *TicketingJobStatusesService
	locales                      *TicketingLocalesService
	macros                       *TicketingMacrosService
	organizationFields           *TicketingOrganizationFieldsService
	organizationMemberships      *TicketingOrganizationMembershipsService
	organizations                *TicketingOrganizationsService
	organizationSubscriptions    *TicketingOrganizationSubscriptionsService
	requests                     *TicketingRequestsService
	resourceCollections          *TicketingResourceCollectionsService
	satisfactionRatings          *TicketingSatisfactionRatingsService
	satisfactionReasons          *TicketingSatisfactionReasonsService
	schedules                    *TicketingSchedulesService
	search                       *TicketingSearchService
	sessions                     *TicketingSessionsService
	sharingAgreements            *TicketingSharingAgreementsService
	sideConversationAttachments  *TicketingSideConversationAttachmentsService
	sideConversationEvents       *TicketingSideConversationEventsService
	sideConversations            *TicketingSideConversationsService
	skillBasedRouting            *TicketingSkillBasedRoutingService
	slaPolicies                  *TicketingSLAPoliciesService
	supportAddresses             *TicketingSupportAddressesService
	suspendedTickets             *TicketingSuspendedTicketsService
	tags                         *TicketingTagsService
	targetFailures               *TicketingTargetFailuresService
	targets                      *TicketingTargetsService
	ticketAttachments            *TicketingTicketAttachmentsService
	ticketAudits                 *TicketingTicketAuditsService
	ticketComments               *TicketingTicketCommentsService
	ticketFields                 *TicketingTicketFieldsService
	ticketForms                  *TicketingTicketFormsService
	ticketImport                 *TicketingTicketImportService
	ticketMetricEvents           *TicketingTicketMetricEventsService
	ticketMetrics                *TicketingTicketMetricsService
	tickets                      *TicketingTicketsService
	ticketSkips                  *TicketingTicketSkipsService
	triggerCategories            *TicketingTriggerCategoriesService
	triggers                     *TicketingTriggersService
	userEvents                   *TicketingUserEventsService
	userFields                   *TicketingUserFieldsService
	userIdentities               *TicketingUserIdentitiesService
	userPasswords                *TicketingUserPasswordsService
	userProfiles                 *TicketingUserProfilesService
	users                        *TicketingUsersService
	views                        *TicketingViewsService
	workspaces                   *TicketingWorkspacesService
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/account_settings/
func (s *TicketingService) AccountSettings() *TicketingAccountSettingsService {
	return s.accountSettings
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/activity_stream/
func (s *TicketingService) ActivityStream() *TicketingActivityStreamService {
	return s.activityStream
}

// https://developer.zendesk.com/api-reference/ticketing/apps/app_location_installations/
func (s *TicketingService) AppLocationInstallations() *TicketingAppLocationInstallationsService {
	return s.appLocationInstallations
}

// https://developer.zendesk.com/api-reference/ticketing/apps/app_locations/
func (s *TicketingService) AppLocations() *TicketingAppLocationsService {
	return s.appLocations
}

// https://developer.zendesk.com/api-reference/ticketing/apps/apps/
func (s *TicketingService) Apps() *TicketingAppsService {
	return s.apps
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/audit_logs/
func (s *TicketingService) AuditLogs() *TicketingAuditLogsService {
	return s.auditLogs
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/
func (s *TicketingService) Automations() *TicketingAutomationsService {
	return s.automations
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/bookmarks/
func (s *TicketingService) Bookmarks() *TicketingBookmarksService {
	return s.bookmarks
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/
func (s *TicketingService) Brands() *TicketingBrandsService {
	return s.brands
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
func (s *TicketingService) CustomRoles() *TicketingCustomRolesService {
	return s.customRoles
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/
func (s *TicketingService) CustomTicketStatuses() *TicketingCustomTicketStatusesService {
	return s.customTicketStatuses
}

// https://developer.zendesk.com/api-reference/ticketing/groups/group_memberships/
func (s *TicketingService) GroupMemberships() *TicketingGroupMembershipsService {
	return s.groupMemberships
}

// https://developer.zendesk.com/api-reference/ticketing/groups/groups/
func (s *TicketingService) Groups() *TicketingGroupsService {
	return s.groups
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/group_sla_policies/
func (s *TicketingService) GroupSLAPolicies() *TicketingGroupSLAPoliciesService {
	return s.groupSLAPolicies
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_exports/
func (s *TicketingService) IncrementalExports() *TicketingIncrementalExportsService {
	return s.incrementalExports
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/incremental_skill_based_routing/
func (s *TicketingService) IncrementalSkillBasedRouting() *TicketingIncrementalSkillBasedRoutingService {
	return s.incrementalSkillBasedRouting
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/job_statuses/
func (s *TicketingService) JobStatuses() *TicketingJobStatusesService {
	return s.jobStatuses
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/
func (s *TicketingService) Locales() *TicketingLocalesService {
	return s.locales
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/macros/
func (s *TicketingService) Macros() *TicketingMacrosService {
	return s.macros
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_fields/
func (s *TicketingService) OrganizationFields() *TicketingOrganizationFieldsService {
	return s.organizationFields
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/
func (s *TicketingService) OrganizationMemberships() *TicketingOrganizationMembershipsService {
	return s.organizationMemberships
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
func (s *TicketingService) Organizations() *TicketingOrganizationsService {
	return s.organizations
}

// https://developer.zendesk.com/api-reference/ticketing/organizations/organization_subscriptions/
func (s *TicketingService) OrganizationSubscriptions() *TicketingOrganizationSubscriptionsService {
	return s.organizationSubscriptions
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-requests/
func (s *TicketingService) Requests() *TicketingRequestsService {
	return s.requests
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/resource_collections/
func (s *TicketingService) ResourceCollections() *TicketingResourceCollectionsService {
	return s.resourceCollections
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/
func (s *TicketingService) SatisfactionRatings() *TicketingSatisfactionRatingsService {
	return s.satisfactionRatings
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_reasons/
func (s *TicketingService) SatisfactionReasons() *TicketingSatisfactionReasonsService {
	return s.satisfactionReasons
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/schedules/
func (s *TicketingService) Schedules() *TicketingSchedulesService {
	return s.schedules
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/search/
func (s *TicketingService) Search() *TicketingSearchService {
	return s.search
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/sessions/
func (s *TicketingService) Sessions() *TicketingSessionsService {
	return s.sessions
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/sharing_agreements/
func (s *TicketingService) SharingAgreements() *TicketingSharingAgreementsService {
	return s.sharingAgreements
}

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation_attachment/
func (s *TicketingService) SideConversationAttachments() *TicketingSideConversationAttachmentsService {
	return s.sideConversationAttachments
}

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation_event/
func (s *TicketingService) SideConversationEvents() *TicketingSideConversationEventsService {
	return s.sideConversationEvents
}

// https://developer.zendesk.com/api-reference/ticketing/side_conversation/side_conversation/
func (s *TicketingService) SideConversations() *TicketingSideConversationsService {
	return s.sideConversations
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/skill_based_routing/
func (s *TicketingService) SkillBasedRouting() *TicketingSkillBasedRoutingService {
	return s.skillBasedRouting
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/sla_policies/
func (s *TicketingService) SLAPolicies() *TicketingSLAPoliciesService {
	return s.slaPolicies
}

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/support_addresses/
func (s *TicketingService) SupportAddresses() *TicketingSupportAddressesService {
	return s.supportAddresses
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/suspended_tickets/
func (s *TicketingService) SuspendedTickets() *TicketingSuspendedTicketsService {
	return s.suspendedTickets
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/
func (s *TicketingService) Tags() *TicketingTagsService {
	return s.tags
}

// https://developer.zendesk.com/api-reference/ticketing/targets/target_failures/
func (s *TicketingService) TargetFailures() *TicketingTargetFailuresService {
	return s.targetFailures
}

// https://developer.zendesk.com/api-reference/ticketing/targets/targets/
func (s *TicketingService) Targets() *TicketingTargetsService {
	return s.targets
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-attachments/
func (s *TicketingService) TicketAttachments() *TicketingTicketAttachmentsService {
	return s.ticketAttachments
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_audits/
func (s *TicketingService) TicketAudits() *TicketingTicketAuditsService {
	return s.ticketAudits
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_comments/
func (s *TicketingService) TicketComments() *TicketingTicketCommentsService {
	return s.ticketComments
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/
func (s *TicketingService) TicketFields() *TicketingTicketFieldsService {
	return s.ticketFields
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_forms/
func (s *TicketingService) TicketForms() *TicketingTicketFormsService {
	return s.ticketForms
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_import/
func (s *TicketingService) TicketImport() *TicketingTicketImportService {
	return s.ticketImport
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_metric_events/
func (s *TicketingService) TicketMetricEvents() *TicketingTicketMetricEventsService {
	return s.ticketMetricEvents
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_metrics/
func (s *TicketingService) TicketMetrics() *TicketingTicketMetricsService {
	return s.ticketMetrics
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func (s *TicketingService) Tickets() *TicketingTicketsService {
	return s.tickets
}

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_skips/
func (s *TicketingService) TicketSkips() *TicketingTicketSkipsService {
	return s.ticketSkips
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/trigger_categories/
func (s *TicketingService) TriggerCategories() *TicketingTriggerCategoriesService {
	return s.triggerCategories
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/triggers/
func (s *TicketingService) Triggers() *TicketingTriggersService {
	return s.triggers
}

// https://developer.zendesk.com/api-reference/ticketing/users/events-api/events-api/
func (s *TicketingService) UserEvents() *TicketingUserEventsService {
	return s.userEvents
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_fields/
func (s *TicketingService) UserFields() *TicketingUserFieldsService {
	return s.userFields
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_identities/
func (s *TicketingService) UserIdentities() *TicketingUserIdentitiesService {
	return s.userIdentities
}

// https://developer.zendesk.com/api-reference/ticketing/users/user_passwords/
func (s *TicketingService) UserPasswords() *TicketingUserPasswordsService {
	return s.userPasswords
}

// https://developer.zendesk.com/api-reference/ticketing/users/profiles_api/profiles_api/
func (s *TicketingService) UserProfiles() *TicketingUserProfilesService {
	return s.userProfiles
}

// https://developer.zendesk.com/api-reference/ticketing/users/users/
func (s *TicketingService) Users() *TicketingUsersService {
	return s.users
}

// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/
func (s *TicketingService) Views() *TicketingViewsService {
	return s.views
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/workspaces/
func (s *TicketingService) Workspaces() *TicketingWorkspacesService {
	return s.workspaces
}
