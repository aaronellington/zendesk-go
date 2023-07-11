package zendesk

import "net/http"

func NewService(
	subDomain string,
	auth authentication,
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
		auth:                 auth,
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
	}
}

type Service struct {
	supportService *SupportService
	guideService   *GuideService
}

// https://developer.zendesk.com/api-reference/ticketing/introduction/
func (s *Service) Support() *SupportService {
	return s.supportService
}

// https://developer.zendesk.com/api-reference/help_center/help-center-api/introduction/
func (s *Service) Guide() *GuideService {
	return s.guideService
}
