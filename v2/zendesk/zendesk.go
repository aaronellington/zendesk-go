package zendesk

import (
	"net/http"
	"time"
)

func New(
	subDomain string,
	authentication Authentication,
	opts ...ConfigOption,
) *Service {
	// Default config
	config := internalConfig{
		userAgent: "aaronellington/zendesk-go",
		timeout:   time.Second * 15,
	}

	// Apply supplied config changes
	for _, opt := range opts {
		opt(&config)
	}

	c := &client{
		subDomain:      subDomain,
		authentication: authentication,
		httpClient: &http.Client{
			Transport: config.roundTripper,
			Timeout:   config.timeout,
		},
		userAgent:            config.userAgent,
		requestPreProcessors: config.requestPreProcessors,
	}

	return &Service{
		subDomain: subDomain,
		webhook:   &WebhookService{},
		ticketing: &TicketingService{
			accountSettings: &TicketingAccountSettingsService{
				c: c,
			},
			activityStream: &TicketingActivityStreamService{
				c: c,
			},
			appLocationInstallations: &TicketingAppLocationInstallationsService{
				c: c,
			},
			appLocations: &TicketingAppLocationsService{
				c: c,
			},
			apps: &TicketingAppsService{
				c: c,
			},
			auditLogs: &TicketingAuditLogsService{
				c: c,
			},
			automations: &TicketingAutomationsService{
				c: c,
			},
			bookmarks: &TicketingBookmarksService{
				c: c,
			},
			brands: &TicketingBrandsService{
				c: c,
			},
			customRoles: &TicketingCustomRolesService{
				c: c,
			},
			customTicketStatuses: &TicketingCustomTicketStatusesService{
				c: c,
			},
			groupMemberships: &TicketingGroupMembershipsService{
				c: c,
			},
			groups: &TicketingGroupsService{
				c: c,
			},
			groupSLAPolicies: &TicketingGroupSLAPoliciesService{
				c: c,
			},
			incrementalExports: &TicketingIncrementalExportsService{
				c: c,
			},
			incrementalSkillBasedRouting: &TicketingIncrementalSkillBasedRoutingService{
				c: c,
			},
			jobStatuses: &TicketingJobStatusesService{
				c: c,
			},
			locales: &TicketingLocalesService{
				c: c,
			},
			macros: &TicketingMacrosService{
				c: c,
			},
			organizationFields: &TicketingOrganizationFieldsService{
				c: c,
			},
			organizationMemberships: &TicketingOrganizationMembershipsService{
				c: c,
			},
			organizations: &TicketingOrganizationsService{
				c: c,
			},
			organizationSubscriptions: &TicketingOrganizationSubscriptionsService{
				c: c,
			},
			requests: &TicketingRequestsService{
				c: c,
			},
			resourceCollections: &TicketingResourceCollectionsService{
				c: c,
			},
			satisfactionRatings: &TicketingSatisfactionRatingsService{
				c: c,
			},
			satisfactionReasons: &TicketingSatisfactionReasonsService{
				c: c,
			},
			schedules: &TicketingSchedulesService{
				c: c,
			},
			search: &TicketingSearchService{
				c: c,
			},
			sessions: &TicketingSessionsService{
				c: c,
			},
			sharingAgreements: &TicketingSharingAgreementsService{
				c: c,
			},
			sideConversationAttachments: &TicketingSideConversationAttachmentsService{
				c: c,
			},
			sideConversationEvents: &TicketingSideConversationEventsService{
				c: c,
			},
			sideConversations: &TicketingSideConversationsService{
				c: c,
			},
			skillBasedRouting: &TicketingSkillBasedRoutingService{
				c: c,
			},
			slaPolicies: &TicketingSLAPoliciesService{
				c: c,
			},
			supportAddresses: &TicketingSupportAddressesService{
				c: c,
			},
			suspendedTickets: &TicketingSuspendedTicketsService{
				c: c,
			},
			tags: &TicketingTagsService{
				c: c,
			},
			targetFailures: &TicketingTargetFailuresService{
				c: c,
			},
			targets: &TicketingTargetsService{
				c: c,
			},
			ticketAttachments: &TicketingTicketAttachmentsService{
				c: c,
			},
			ticketAudits: &TicketingTicketAuditsService{
				c: c,
			},
			ticketComments: &TicketingTicketCommentsService{
				c: c,
			},
			ticketFields: &TicketingTicketFieldsService{
				c: c,
			},
			ticketForms: &TicketingTicketFormsService{
				c: c,
			},
			ticketImport: &TicketingTicketImportService{
				c: c,
			},
			ticketMetricEvents: &TicketingTicketMetricEventsService{
				c: c,
			},
			ticketMetrics: &TicketingTicketMetricsService{
				c: c,
			},
			tickets: &TicketingTicketsService{
				c: c,
			},
			ticketSkips: &TicketingTicketSkipsService{
				c: c,
			},
			triggerCategories: &TicketingTriggerCategoriesService{
				c: c,
			},
			triggers: &TicketingTriggersService{
				c: c,
			},
			userEvents: &TicketingUserEventsService{
				c: c,
			},
			userFields: &TicketingUserFieldsService{
				c: c,
			},
			userIdentities: &TicketingUserIdentitiesService{
				c: c,
			},
			userPasswords: &TicketingUserPasswordsService{
				c: c,
			},
			userProfiles: &TicketingUserProfilesService{
				c: c,
			},
			users: &TicketingUsersService{
				c: c,
			},
			views: &TicketingViewsService{
				c: c,
			},
			workspaces: &TicketingWorkspacesService{
				c: c,
			},
		},
		helpCenter: &HelpCenterService{
			accountCustomClaims: &HelpCenterAccountCustomClaimsService{
				c: c,
			},
			articleAttachments: &HelpCenterArticleAttachmentsService{
				c: c,
			},
			articleComments: &HelpCenterArticleCommentsService{
				c: c,
			},
			articleLabels: &HelpCenterArticleLabelsService{
				c: c,
			},
			articles: &HelpCenterArticlesService{
				c: c,
			},
			badgeAssignments: &HelpCenterBadgeAssignmentsService{
				c: c,
			},
			badgeCategories: &HelpCenterBadgeCategoriesService{
				c: c,
			},
			badges: &HelpCenterBadgesService{
				c: c,
			},
			categories: &HelpCenterCategoriesService{
				c: c,
			},
			contentSubscriptions: &HelpCenterContentSubscriptionsService{
				c: c,
			},
			contentTags: &HelpCenterContentTagsService{
				c: c,
			},
			jwts: &HelpCenterJWTsService{
				c: c,
			},
			permissionGroups: &HelpCenterPermissionGroupsService{
				c: c,
			},
			postComments: &HelpCenterPostCommentsService{
				c: c,
			},
			posts: &HelpCenterPostsService{
				c: c,
			},
			search: &HelpCenterSearchService{
				c: c,
			},
			sections: &HelpCenterSectionsService{
				c: c,
			},
			theming: &HelpCenterThemingService{
				c: c,
			},
			topics: &HelpCenterTopicsService{
				c: c,
			},
			translations: &HelpCenterTranslationsService{
				c: c,
			},
			userImages: &HelpCenterUserImagesService{
				c: c,
			},
			userSegments: &HelpCenterUserSegmentsService{
				c: c,
			},
			userSubscriptions: &HelpCenterUserSubscriptionsService{
				c: c,
			},
			votes: &HelpCenterVotesService{
				c: c,
			},
		},
		liveChat: &LiveChatService{
			accounts: &LiveChatAccountsService{
				c: c,
			},
			agents: &LiveChatAgentsService{
				c: c,
			},
			bans: &LiveChatBansService{
				c: c,
			},
			chats: &LiveChatChatsService{
				c: c,
			},
			departments: &LiveChatDepartmentsService{
				c: c,
			},
			goals: &LiveChatGoalsService{
				c: c,
			},
			incrementalAgentEvents: &LiveChatIncrementalAgentEventsService{
				c: c,
			},
			incrementalExports: &LiveChatIncrementalExportsService{
				c: c,
			},
			oauthClients: &LiveChatOAuthClientsService{
				c: c,
			},
			oauthTokens: &LiveChatOAuthTokensService{
				c: c,
			},
			realTime: &LiveChatRealTimeService{
				restService: &LiveChatRealTimeRESTService{
					c: c,
				},
				streamingService: &LiveChatRealTimeStreamingService{
					c: c,
				},
			},
			roles: &LiveChatRolesService{
				c: c,
			},
			routingSettings: &LiveChatRoutingSettingsService{
				c: c,
			},
			shortcuts: &LiveChatShortcutsService{
				c: c,
			},
			skills: &LiveChatSkillsService{
				c: c,
			},
			triggers: &LiveChatTriggersService{
				c: c,
			},
			visitors: &LiveChatVisitorsService{
				c: c,
			},
		},
	}
}

type Service struct {
	helpCenter *HelpCenterService
	liveChat   *LiveChatService
	subDomain  string
	ticketing  *TicketingService
	webhook    *WebhookService
}

func (s *Service) HelpCenter() *HelpCenterService {
	return s.helpCenter
}

func (s *Service) LiveChat() *LiveChatService {
	return s.liveChat
}

func (s *Service) SubDomain() string {
	return s.subDomain
}

func (s *Service) Ticketing() *TicketingService {
	return s.ticketing
}

func (s *Service) Webhook() *WebhookService {
	return s.webhook
}
