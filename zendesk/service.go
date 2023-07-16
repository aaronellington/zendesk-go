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
		subdomain:            subDomain,
		zendeskAuth:          zendeskAuth,
		chatCredentials:      chatCredentials,
		chatMutex:            &sync.Mutex{},
		chatToken:            nil,
		requestPreProcessors: config.requestPreProcessors,
	}

	return &Service{
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
		},
		guideService: &GuideService{
			categoriesService: &CategoriesService{
				client: c,
			},
			sectionsService: &SectionsService{
				client: c,
			},
			articlesService: &ArticlesService{
				client: c,
			},
		},
		chatService: &ChatService{
			chatsService: &ChatsService{
				client: c,
			},
			agentsService: &AgentEventService{
				client:           c,
				agentStatesMutex: &sync.Mutex{},
				agentStates:      AgentStates{},
			},
		},
	}
}

type Service struct {
	supportService *SupportService
	guideService   *GuideService
	chatService    *ChatService
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
func (s *Service) Chat() *ChatService {
	return s.chatService
}
