package zendesk

import (
	"net/http"
	"sync"
	"time"
)

func NewService(
	subDomain string,
	zendeskAuth authentication,
	chatCredentials ChatCredentials,
	opts ...configOption,
) *Service {
	config := &internalConfig{
		userAgent: "aaronellington/zendesk-go",
		timeout:   time.Second * 15,
	}
	for _, opt := range opts {
		opt(config)
	}

	c := &client{
		httpClientForZendesk: &http.Client{
			Transport: config.roundTripper,
			Timeout:   config.timeout,
		},
		httpClientForZopim: &http.Client{
			Transport: config.roundTripper,
			Timeout:   config.timeout,
		},
		userAgent:            config.userAgent,
		subDomain:            subDomain,
		zendeskAuth:          zendeskAuth,
		chatCredentials:      chatCredentials,
		chatMutex:            &sync.Mutex{},
		chatToken:            nil,
		requestPreProcessors: config.requestPreProcessors,
	}

	return &Service{
		subDomain: subDomain,
		accountConfigurationService: &AccountConfigurationService{
			auditLogService: &AuditLogService{
				client: c,
				generic: genericService[
					AuditLogID,
					AuditLogResponse,
					AuditLogsResponse,
				]{
					client:  c,
					apiName: "audit_logs",
				},
			},
			brandService: &BrandService{
				client: c,
				generic: genericService[
					BrandID,
					BrandResponse,
					BrandsResponse,
				]{
					client:  c,
					apiName: "brands",
				},
			},
			customRoleService: &CustomRoleService{
				client: c,
				generic: genericService[
					CustomRoleID,
					CustomRoleResponse,
					CustomRolesResponse,
				]{
					client:  c,
					apiName: "custom_roles",
				},
			},
		},
		supportService: &SupportService{
			customStatusService: &CustomStatusService{
				client: c,
				generic: genericService[
					CustomStatusID,
					CustomStatusResponse,
					CustomStatusesResponse,
				]{
					client:  c,
					apiName: "custom_statuses",
				},
			},
			organizationService: &OrganizationService{
				client: c,
				generic: genericService[
					OrganizationID,
					OrganizationResponse,
					OrganizationsResponse,
				]{
					client:  c,
					apiName: "organizations",
				},
			},
			ticketService: &TicketService{
				client: c,
				generic: genericService[
					TicketID,
					TicketResponse,
					TicketsResponse,
				]{
					client:  c,
					apiName: "tickets",
				},
			},
			ticketAuditService: &TicketAuditService{
				client: c,
			},
			sideConversationService: &SideConversationService{
				client: c,
			},
			scheduleService: &ScheduleService{
				client: c,
			},
			groupMembershipService: &GroupMembershipService{
				client: c,
				generic: genericService[
					GroupMembershipID,
					GroupMembershipResponse,
					GroupMembershipsResponse,
				]{
					client:  c,
					apiName: "group_memberships",
				},
			},
			groupService: &GroupsService{
				client: c,
				generic: genericService[
					GroupID,
					GroupResponse,
					GroupsResponse,
				]{
					client:  c,
					apiName: "groups",
				},
			},
			organizationFieldService: &OrganizationFieldService{
				client: c,
				generic: genericService[
					OrganizationFieldID,
					OrganizationFieldResponse,
					OrganizationFieldsResponse,
				]{
					client:  c,
					apiName: "organization_fields",
				},
			},
			organizationMembershipService: &OrganizationMembershipService{
				client: c,
				generic: genericService[
					OrganizationMembershipID,
					OrganizationMembershipResponse,
					OrganizationMembershipsResponse,
				]{
					client:  c,
					apiName: "organization_memberships",
				},
			},
			suspendedTicketService: &SuspendedTicketService{
				client: c,
				generic: genericService[
					SuspendedTicketID,
					SuspendedTicketResponse,
					SuspendedTicketsResponse,
				]{
					client:  c,
					apiName: "suspended_tickets",
				},
			},
			satisfactionRatingService: &SatisfactionRatingService{
				client: c,
				generic: genericService[
					SatisfactionRatingID,
					SatisfactionRatingResponse,
					SatisfactionRatingsResponse,
				]{
					client:  c,
					apiName: "satisfaction_ratings",
				},
			},
			ticketAttachmentService: &TicketAttachmentService{
				client: c,
			},
			ticketCommentService: &TicketCommentService{
				client: c,
			},
			viewService: &ViewService{
				client: c,
				generic: genericService[
					ViewID,
					ViewResponse,
					ViewsResponse,
				]{
					client:  c,
					apiName: "views",
				},
			},
			ticketFormService: &TicketFormService{
				client: c,
				generic: genericService[
					TicketFormID,
					TicketFormResponse,
					TicketFormsResponse,
				]{
					client:  c,
					apiName: "ticket_forms",
				},
			},
			ticketFieldService: &TicketFieldService{
				client: c,
				generic: genericService[
					TicketFieldID,
					TicketFieldResponse,
					TicketFieldsResponse,
				]{
					client:  c,
					apiName: "ticket_fields",
				},
			},
			userFieldsService: &UserFieldService{
				client: c,
				generic: genericService[
					UserFieldID,
					UserFieldResponse,
					UserFieldsResponse,
				]{
					client:  c,
					apiName: "user_fields",
				},
			},
			ticketTagService: &TicketTagService{
				client: c,
			},
			macroService: &MacroService{
				client: c,
				generic: genericService[
					MacroID,
					MacroResponse,
					MacrosResponse,
				]{
					client:  c,
					apiName: "macros",
				},
			},
			automationService: &AutomationService{
				client: c,
				generic: genericService[
					AutomationID,
					AutomationResponse,
					AutomationsResponse,
				]{
					client:  c,
					apiName: "automations",
				},
			},
			triggerService: &TriggerService{
				client: c,
				generic: genericService[
					TriggerID,
					TriggerResponse,
					TriggersResponse,
				]{
					client:  c,
					apiName: "triggers",
				},
			},
			userIdentityService: &UserIdentityService{
				client: c,
			},
			userService: &UserService{
				client: c,
				generic: genericService[
					UserID,
					UserResponse,
					UsersResponse,
				]{
					client:  c,
					apiName: "users",
				},
			},
		},
		guideService: &GuideService{
			categoriesService: &CategoryService{
				client: c,
				generic: genericService[
					CategoryID,
					CategoryResponse,
					CategoriesResponse,
				]{
					client:  c,
					apiName: "help_center/categories",
				},
			},
			sectionsService: &SectionService{
				client: c,
				generic: genericService[
					SectionID,
					SectionResponse,
					SectionsResponse,
				]{
					client:  c,
					apiName: "help_center/sections",
				},
			},
			articlesService: &ArticleService{
				client: c,
				generic: genericService[
					ArticleID,
					ArticleResponse,
					ArticlesResponse,
				]{
					client:  c,
					apiName: "help_center/articles",
				},
			},
		},
		liveChatService: &LiveChatService{
			chatService: &ChatService{
				chatsService: &ChatsService{
					client: c,
				},
				oauthClientService: &OAuthClientService{
					client: c,
				},
				agentEventService: &AgentEventService{
					client:           c,
					agentStatesMutex: &sync.Mutex{},
					agentStates:      AgentStates{},
				},
				departmentService: &DepartmentService{
					client: c,
				},
			},
			realTimeChatService: &RealTimeChatService{
				realTimeChatRestService: &RealTimeChatRestService{
					client: c,
				},
				realTimeChatStreamingService: &RealTimeChatStreamingService{
					client: c,
				},
			},
			chatConversationsService: &ChatConversationsService{
				client: c,
			},
		},
		webhookService: &WebhookService{
			client: c,
		},
	}
}

type Service struct {
	subDomain                   string
	accountConfigurationService *AccountConfigurationService
	supportService              *SupportService
	guideService                *GuideService
	liveChatService             *LiveChatService
	webhookService              *WebhookService
}

func (s *Service) SubDomain() string {
	return s.subDomain
}

// https://developer.zendesk.com/api-reference/ticketing/introduction/
func (s *Service) Support() *SupportService {
	return s.supportService
}

// https://developer.zendesk.com/api-reference/ticketing/introduction/
func (s *Service) AccountConfiguration() *AccountConfigurationService {
	return s.accountConfigurationService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
func (s *Service) Guide() *GuideService {
	return s.guideService
}

// https://developer.zendesk.com/api-reference/live-chat/introduction/
func (s *Service) LiveChat() *LiveChatService {
	return s.liveChatService
}

// https://developer.zendesk.com/api-reference/webhooks/introduction/
func (s *Service) Webhook() *WebhookService {
	return s.webhookService
}
