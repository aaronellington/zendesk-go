package zendesk

import (
	"net/http"
	"sync"
)

func NewService(
	subDomain string,
	zendeskAuth authentication,
	chatCredentials ChatCredentials,
	opts ...configOption,
) *Service {
	config := &internalConfig{}
	for _, opt := range opts {
		opt(config)
	}

	c := &client{
		httpClient: &http.Client{
			Transport: config.roundTripper,
		},
		subDomain:            subDomain,
		zendeskAuth:          zendeskAuth,
		chatCredentials:      chatCredentials,
		chatMutex:            &sync.Mutex{},
		chatToken:            nil,
		requestPreProcessors: config.requestPreProcessors,
	}

	return &Service{
		subDomain: subDomain,
		supportService: &SupportService{
			organizationService: &OrganizationService{
				client: c,
			},
			userService: &UserService{
				client: c,
			},
			ticketService: &TicketService{
				client: c,
			},
			ticketAuditService: &TicketAuditService{
				client: c,
			},
			ticketMetricEventService: &TicketMetricEventService{
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
			suspendedTicketService: &SuspendedTicketService{
				client: c,
			},
			userFieldsService: &UserFieldService{
				client: c,
			},
			userIdentityService: &UserIdentityService{
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
				client: c,
			},
			agentEventService: &AgentEventService{
				client:           c,
				agentStatesMutex: &sync.Mutex{},
				agentStates:      AgentStates{},
			},
		},
	}
}

type Service struct {
	subDomain       string
	supportService  *SupportService
	guideService    *GuideService
	liveChatService *LiveChatService
}

func (s *Service) SubDomain() string {
	return s.subDomain
}

// https://developer.zendesk.com/api-reference/ticketing/introduction/
func (s *Service) Support() *SupportService {
	return s.supportService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
func (s *Service) Guide() *GuideService {
	return s.guideService
}

// https://developer.zendesk.com/api-reference/live-chat/introduction/
func (s *Service) LiveChat() *LiveChatService {
	return s.liveChatService
}
