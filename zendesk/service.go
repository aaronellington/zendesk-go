package zendesk

import (
	"net/http"
	"sync"
	"time"

	"github.com/aaronellington/zendesk-go/zendesk/internal/utils"
)

func NewService(
	subDomain string,
	zendeskAuth authentication,
	chatCredentials ChatCredentials,
	opts ...ConfigOption,
) *Service {
	config := &internalConfig{
		userAgent:                 "aaronellington/zendesk-go",
		timeout:                   time.Second * 15,
		realTimeChatWebsocketHost: RealTimeChatStreamingHost,
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

	wsClient := wsClient{
		client: c,
		conn:   nil,
		// connEstablished: make(chan bool, 1),
		connMutex: &sync.Mutex{},
		rtcWSHost: config.realTimeChatWebsocketHost,
	}

	return &Service{
		subDomain: subDomain,
		accountConfigurationService: &AccountConfigurationService{
			auditLogService: &AuditLogService{
				client: c,
			},
			brandService: &BrandService{
				client: c,
			},
			customRoleService: &CustomRoleService{
				client: c,
			},
		},
		supportService: &SupportService{
			customStatusService: &CustomStatusService{
				client: c,
			},
			organizationService: &OrganizationService{
				client: c,
			},
			ticketService: &TicketService{
				client: c,
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
			},
			groupService: &GroupsService{
				client: c,
			},
			organizationFieldService: &OrganizationFieldService{
				client: c,
			},
			organizationMembershipService: &OrganizationMembershipService{
				client: c,
			},
			suspendedTicketService: &SuspendedTicketService{
				client: c,
			},
			satisfactionRatingService: &SatisfactionRatingService{
				client: c,
			},
			ticketAttachmentService: &TicketAttachmentService{
				client: c,
			},
			ticketCommentService: &TicketCommentService{
				client: c,
			},
			viewService: &ViewService{
				client: c,
			},
			ticketFormService: &TicketFormService{
				client: c,
			},
			ticketFieldService: &TicketFieldService{
				client: c,
			},
			userFieldsService: &UserFieldService{
				client: c,
			},
			ticketTagService: &TicketTagService{
				client: c,
			},
			macroService: &MacroService{
				client: c,
			},
			automationService: &AutomationService{
				client: c,
			},
			triggerService: &TriggerService{
				client: c,
			},
			userIdentityService: &UserIdentityService{
				client: c,
			},
			userService: &UserService{
				client: c,
			},
		},
		guideService: &GuideService{
			categoriesService: &CategoryService{
				client: c,
			},
			sectionsService: &SectionService{
				client: c,
			},
			articlesService: &ArticleService{
				client: c,
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
					wsClient: &wsClient,
					wsCache: &wsCache{
						chat:  utils.NewMemoryCacheInstance[GroupID, WebsocketChatMetricData](),
						agent: utils.NewMemoryCacheInstance[GroupID, WebsocketAgentMetricData](),
						metadata: &wsConnMetadata{
							mutex: &sync.Mutex{},
						},
					},
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
